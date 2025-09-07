package services

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
	"github.com/farislr/daneizo/internal/loanv2/domain/repositories"
)

// FraudRiskLevel represents different fraud risk levels
type FraudRiskLevel int

const (
	NoFraudRisk FraudRiskLevel = iota
	LowFraudRisk
	MediumFraudRisk
	HighFraudRisk
)

// RiskAssessmentResult contains the comprehensive risk assessment
type RiskAssessmentResult struct {
	OverallRisk       valueobjects.RiskCategory
	CreditScore       valueobjects.CreditScore
	FraudRisk         FraudRiskLevel
	DebtToIncomeRatio decimal.Decimal
	RecommendedRate   valueobjects.InterestRate
	MaxLoanAmount     valueobjects.Money
	Reasoning         []string
}

// External service ports for risk assessment
type CreditBureauPort interface {
	GetCreditScore(ctx context.Context, userID valueobjects.UserID) (valueobjects.CreditScore, error)
	GetCreditHistory(ctx context.Context, userID valueobjects.UserID) (CreditHistory, error)
}

type FraudDetectionPort interface {
	AssessFraudRisk(ctx context.Context, userID valueobjects.UserID) (FraudRiskLevel, error)
	ReportSuspiciousActivity(ctx context.Context, userID valueobjects.UserID, activity string) error
}

type CreditHistory struct {
	NumberOfAccounts    int
	PaymentHistory      []PaymentRecord
	DelinquencyHistory  []DelinquencyRecord
	TotalDebt          valueobjects.Money
	CreditUtilization  decimal.Decimal
}

type PaymentRecord struct {
	Date   string
	Amount valueobjects.Money
	OnTime bool
}

type DelinquencyRecord struct {
	Date         string
	Amount       valueobjects.Money
	DaysOverdue  int
	Resolved     bool
}

// RiskAssessmentService performs comprehensive credit and fraud risk assessment
type RiskAssessmentService struct {
	creditBureauPort   CreditBureauPort
	fraudDetectionPort FraudDetectionPort
	userRepository     repositories.UserRepository
}

// NewRiskAssessmentService creates a new risk assessment service
func NewRiskAssessmentService(
	creditBureauPort CreditBureauPort,
	fraudDetectionPort FraudDetectionPort,
	userRepository repositories.UserRepository,
) *RiskAssessmentService {
	return &RiskAssessmentService{
		creditBureauPort:   creditBureauPort,
		fraudDetectionPort: fraudDetectionPort,
		userRepository:     userRepository,
	}
}

// AssessLoanRisk performs comprehensive risk assessment for a loan application
func (ras *RiskAssessmentService) AssessLoanRisk(
	ctx context.Context,
	borrowerID valueobjects.UserID,
	loanAmount valueobjects.Money,
	termMonths int,
) (*RiskAssessmentResult, error) {
	result := &RiskAssessmentResult{
		Reasoning: make([]string, 0),
	}

	// 1. Get borrower information
	borrower, err := ras.userRepository.FindByID(ctx, borrowerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get borrower: %w", err)
	}

	// 2. Get credit score from external credit bureau
	creditScore, err := ras.creditBureauPort.GetCreditScore(ctx, borrowerID)
	if err != nil {
		return nil, fmt.Errorf("credit assessment failed: %w", err)
	}
	result.CreditScore = creditScore

	// 3. Perform fraud detection check
	fraudRisk, err := ras.fraudDetectionPort.AssessFraudRisk(ctx, borrowerID)
	if err != nil {
		// Don't fail the assessment, but log the issue
		fraudRisk = MediumFraudRisk
		result.Reasoning = append(result.Reasoning, "Could not complete fraud assessment")
	}
	result.FraudRisk = fraudRisk

	// 4. Calculate debt-to-income ratio
	debtToIncomeRatio := ras.calculateDebtToIncomeRatio(borrower, loanAmount, termMonths)
	result.DebtToIncomeRatio = debtToIncomeRatio

	// 5. Determine overall risk category using business rules
	overallRisk := ras.determineOverallRisk(creditScore, fraudRisk, debtToIncomeRatio, loanAmount)
	result.OverallRisk = overallRisk

	// 6. Calculate recommended interest rate based on risk
	recommendedRate, err := ras.calculateRecommendedRate(overallRisk, creditScore, loanAmount, termMonths)
	if err != nil {
		return nil, fmt.Errorf("rate calculation failed: %w", err)
	}
	result.RecommendedRate = recommendedRate

	// 7. Determine maximum loan amount for this borrower
	maxLoanAmount := ras.calculateMaxLoanAmount(creditScore, borrower.MonthlyIncome, debtToIncomeRatio)
	result.MaxLoanAmount = maxLoanAmount

	// 8. Add business reasoning
	result.Reasoning = append(result.Reasoning, ras.generateReasoning(result)...)

	return result, nil
}

// determineOverallRisk applies business rules to determine overall risk
func (ras *RiskAssessmentService) determineOverallRisk(
	creditScore valueobjects.CreditScore,
	fraudRisk FraudRiskLevel,
	debtToIncomeRatio decimal.Decimal,
	loanAmount valueobjects.Money,
) valueobjects.RiskCategory {
	// Immediate disqualifiers
	if fraudRisk == HighFraudRisk {
		return valueobjects.UnacceptableRisk
	}

	if debtToIncomeRatio.GreaterThan(decimal.NewFromFloat(0.5)) { // >50% DTI
		return valueobjects.UnacceptableRisk
	}

	// Excellent credit with low risk factors
	if creditScore.IsExcellent() && fraudRisk == NoFraudRisk &&
		debtToIncomeRatio.LessThan(decimal.NewFromFloat(0.3)) &&
		loanAmount.Amount().LessThan(decimal.NewFromInt(100000)) {
		return valueobjects.LowRisk
	}

	// Good credit with reasonable risk factors
	if creditScore.IsGood() && fraudRisk <= LowFraudRisk &&
		debtToIncomeRatio.LessThan(decimal.NewFromFloat(0.4)) {
		return valueobjects.MediumRisk
	}

	// Fair credit or higher risk factors
	if creditScore.IsFair() || fraudRisk == MediumFraudRisk ||
		debtToIncomeRatio.GreaterThan(decimal.NewFromFloat(0.4)) {
		return valueobjects.HighRisk
	}

	// Default to high risk for poor credit
	return valueobjects.HighRisk
}

// calculateRecommendedRate calculates interest rate based on risk factors
func (ras *RiskAssessmentService) calculateRecommendedRate(
	riskCategory valueobjects.RiskCategory,
	creditScore valueobjects.CreditScore,
	loanAmount valueobjects.Money,
	termMonths int,
) (valueobjects.InterestRate, error) {
	baseRate := decimal.NewFromFloat(8.0) // 8% base rate

	// Risk-based adjustments
	switch riskCategory {
	case valueobjects.LowRisk:
		baseRate = baseRate.Sub(decimal.NewFromFloat(2.0)) // -2% for low risk
	case valueobjects.MediumRisk:
		// No adjustment
	case valueobjects.HighRisk:
		baseRate = baseRate.Add(decimal.NewFromFloat(3.0)) // +3% for high risk
	case valueobjects.UnacceptableRisk:
		return valueobjects.InterestRate{}, fmt.Errorf("loan cannot be approved due to unacceptable risk")
	}

	// Credit score fine-tuning
	if creditScore.IsExcellent() {
		baseRate = baseRate.Sub(decimal.NewFromFloat(0.5)) // Additional -0.5% for excellent credit
	} else if creditScore.IsPoor() {
		baseRate = baseRate.Add(decimal.NewFromFloat(1.0)) // Additional +1% for poor credit
	}

	// Loan amount adjustments (smaller loans = higher rates due to fixed costs)
	if loanAmount.Amount().LessThan(decimal.NewFromInt(10000)) {
		baseRate = baseRate.Add(decimal.NewFromFloat(1.5)) // +1.5% for small loans
	} else if loanAmount.Amount().GreaterThan(decimal.NewFromInt(100000)) {
		baseRate = baseRate.Sub(decimal.NewFromFloat(0.5)) // -0.5% for large loans
	}

	// Term adjustments (longer terms = higher rates)
	if termMonths > 60 {
		baseRate = baseRate.Add(decimal.NewFromFloat(1.0)) // +1% for long terms
	} else if termMonths <= 12 {
		baseRate = baseRate.Sub(decimal.NewFromFloat(0.5)) // -0.5% for short terms
	}

	// Ensure rate stays within acceptable bounds
	minRate := decimal.NewFromFloat(5.0)
	maxRate := decimal.NewFromFloat(35.0)

	if baseRate.LessThan(minRate) {
		baseRate = minRate
	}
	if baseRate.GreaterThan(maxRate) {
		baseRate = maxRate
	}

	return valueobjects.NewInterestRate(baseRate)
}

// calculateDebtToIncomeRatio calculates the debt-to-income ratio including the new loan
func (ras *RiskAssessmentService) calculateDebtToIncomeRatio(
	borrower *repositories.User,
	loanAmount valueobjects.Money,
	termMonths int,
) decimal.Decimal {
	if borrower.MonthlyIncome.IsZero() {
		return decimal.NewFromFloat(1.0) // 100% DTI if no income reported
	}

	// Calculate estimated monthly payment (simplified calculation)
	interestRate := decimal.NewFromFloat(10.0) // Assume 10% for DTI calculation
	monthlyRate := interestRate.Div(decimal.NewFromFloat(100)).Div(decimal.NewFromFloat(12))

	// Monthly payment formula: P * [r(1+r)^n] / [(1+r)^n - 1]
	onePlusR := decimal.NewFromFloat(1).Add(monthlyRate)
	onePlusRPowN := onePlusR.Pow(decimal.NewFromInt(int64(termMonths)))

	monthlyPayment := loanAmount.Amount().Mul(
		monthlyRate.Mul(onePlusRPowN).Div(onePlusRPowN.Sub(decimal.NewFromFloat(1))),
	)

	totalMonthlyDebt := borrower.MonthlyDebtPayment.Amount().Add(monthlyPayment)

	return totalMonthlyDebt.Div(borrower.MonthlyIncome.Amount())
}

// calculateMaxLoanAmount calculates maximum borrowing capacity
func (ras *RiskAssessmentService) calculateMaxLoanAmount(
	creditScore valueobjects.CreditScore,
	income valueobjects.Money,
	currentDTI decimal.Decimal,
) valueobjects.Money {
	maxDTI := decimal.NewFromFloat(0.43) // 43% maximum DTI

	if currentDTI.GreaterThanOrEqual(maxDTI) {
		return valueobjects.NewZeroMoney(income.Currency()) // No additional borrowing capacity
	}

	availableDTI := maxDTI.Sub(currentDTI)
	availablePayment := income.Amount().Mul(availableDTI)

	// Convert available payment to loan amount (simplified)
	// Assume 10% rate, 36 month term for calculation
	estimatedLoanAmount := availablePayment.Mul(decimal.NewFromInt(30)) // Rough multiplier

	// Apply credit score limits
	var maxAmount decimal.Decimal
	switch {
	case creditScore.IsExcellent():
		maxAmount = decimal.NewFromInt(250000)
	case creditScore.IsGood():
		maxAmount = decimal.NewFromInt(150000)
	case creditScore.IsFair():
		maxAmount = decimal.NewFromInt(75000)
	default:
		maxAmount = decimal.NewFromInt(25000)
	}

	if estimatedLoanAmount.GreaterThan(maxAmount) {
		estimatedLoanAmount = maxAmount
	}

	result, _ := valueobjects.NewMoney(estimatedLoanAmount, income.Currency())
	return result
}

// generateReasoning creates human-readable reasoning for the risk assessment
func (ras *RiskAssessmentService) generateReasoning(result *RiskAssessmentResult) []string {
	reasoning := make([]string, 0)

	// Credit score reasoning
	if result.CreditScore.IsExcellent() {
		reasoning = append(reasoning, "Excellent credit score indicates strong payment history")
	} else if result.CreditScore.IsPoor() {
		reasoning = append(reasoning, "Poor credit score indicates payment difficulties")
	}

	// Fraud risk reasoning
	if result.FraudRisk == HighFraudRisk {
		reasoning = append(reasoning, "High fraud risk detected - loan denied")
	} else if result.FraudRisk == NoFraudRisk {
		reasoning = append(reasoning, "No fraud indicators detected")
	}

	// DTI reasoning
	if result.DebtToIncomeRatio.GreaterThan(decimal.NewFromFloat(0.4)) {
		reasoning = append(reasoning, "High debt-to-income ratio indicates potential payment stress")
	} else if result.DebtToIncomeRatio.LessThan(decimal.NewFromFloat(0.3)) {
		reasoning = append(reasoning, "Low debt-to-income ratio indicates good payment capacity")
	}

	// Overall risk reasoning
	switch result.OverallRisk {
	case valueobjects.LowRisk:
		reasoning = append(reasoning, "Low risk profile qualifies for preferential rates")
	case valueobjects.HighRisk:
		reasoning = append(reasoning, "High risk profile requires elevated interest rate")
	case valueobjects.UnacceptableRisk:
		reasoning = append(reasoning, "Risk profile exceeds lending criteria")
	}

	return reasoning
}