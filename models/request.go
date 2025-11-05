package models

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ConvertPDF struct {
	TransactionID       string `json:"transactionId"`
	Status              string `json:"status"`
	DebitAccountName    string `json:"debitAccountName"`
	DebitAccountNumber  string `json:"debitAccountNumber"`
	ModuleType          string `json:"moduleType"`
	TotalPayment        string `json:"totalPayment"`
	NetValue            string `json:"netValue"`
	PPNTax              string `json:"ppnTax"`
	PBBKBTax            string `json:"pbbkbTax"`
	PPHTax              string `json:"pphTax"`
	GrossValue          string `json:"grossValue"`
	DebitCreditValue    string `json:"debitCreditValue"`
	TotalAmount         string `json:"totalAmount"`
	SoldToName          string `json:"soldToName"`
	Buyer               string `json:"buyer"`
	DepoName            string `json:"depoName"`
	SalesOrganization   string `json:"salesOrganization"`
	ProductGroup        string `json:"productGroup"`
	DistributionChannel string `json:"distributionChannel"`
	Payer               string `json:"payer"`
}

type GeneratePDFRequest struct {
	ProcessID     string `json:"process_id"`
	TransactionID string `json:"transaction_id"`
	Detail        Detail `json:"detail"`
	Logs          []Logs `json:"logs"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"` // base64 PDF
	Error   string `json:"error,omitempty"`
}

type Logs struct {
	Command     string    `json:"command"`
	Type        string    `json:"type"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	Username    string    `json:"username"`
	CompanyName string    `json:"companyName"`
	Task        any       `json:"task"`
	CreatedAt   time.Time `json:"createdAt"`
	TaskID      string    `json:"taskID"`
	Key         string    `json:"key"`
}

type Detail struct {
	TransactionInformation struct {
		ApplicationID             string `json:"applicationId"`
		DebitAccountNumber        string `json:"debitAccountNumber"`
		DebitAccountName          string `json:"debitAccountName"`
		TransactionDate           string `json:"transactionDate"`
		ShipTo                    string `json:"shipTo"`
		ProductGroupID            string `json:"productGroupId"`
		ProductGroupName          string `json:"productGroupName"`
		Status                    string `json:"status"`
		QuotationNumber           string `json:"quotationNumber"`
		SchedulingAgreementNumber string `json:"schedulingAgreementNumber"`
		ErrorMessagePayment       string `json:"errorMessagePayment"`
		DepoID                    string `json:"depoId"`
		DepoName                  string `json:"depoName"`
		PayerID                   string `json:"payerId"`
		PayerName                 string `json:"payerName"`
		ShipToDesc                string `json:"shipToDesc"`
		SoldTo                    string `json:"soldTo"`
		SoldToDesc                string `json:"soldToDesc"`
	} `json:"transactionInformation"`
	PurchaseDetail struct {
		NetValue         int `json:"netValue"`
		PpnTax           int `json:"ppnTax"`
		PbbkbTax         int `json:"pbbkbTax"`
		PphTax           int `json:"pphTax"`
		GrossValue       int `json:"grossValue"`
		Fee              int `json:"fee"`
		DebitCreditValue int `json:"debitCreditValue"`
		TotalAmount      int `json:"totalAmount"`
		TotalPayment     int `json:"totalPayment"`
	} `json:"purchaseDetail"`
	Task struct {
		TaskID             string `json:"taskId"`
		Status             string `json:"status"`
		Reasons            string `json:"reasons"`
		LastApprovedByName string `json:"lastApprovedByName"`
		LastRejectedByName string `json:"lastRejectedByName"`
		UpdatedByName      string `json:"updatedByName"`
		CreatedByName      string `json:"createdByName"`
		Type               string `json:"type"`
		Comment            string `json:"comment"`
		CreatedAt          string `json:"createdAt"`
		UpdatedAt          string `json:"updatedAt"`
		ApprovedAt         string `json:"approvedAt"`
		RejectedAt         string `json:"rejectedAt"`
	} `json:"task"`
	WorkflowDoc         string `json:"workflowDoc"`
	MaterialInformation []struct {
		MaterialID          string `json:"materialId"`
		MaterialDescription string `json:"materialDescription"`
		DeliveryDate        string `json:"deliveryDate"`
		Transporter         string `json:"transporter"`
		Trip                string `json:"trip"`
		Quantity            int    `json:"quantity"`
		Uom                 string `json:"uom"`
	} `json:"materialInformation"`
}

type Workflow struct {
	Header         *Header                `json:"header"`
	Records        *Records               `json:"records"`
	CreatedBy      *User                  `json:"createdBy"`
	CreatedAt      *timestamppb.Timestamp `json:"createdAt"`
	CurrentRoleIDs []int                  `json:"currentRoleIDs"`
	CurrentStep    string                 `json:"currentStep"`
}

type Header struct {
	ProductID           int    `json:"productID"`
	ProductName         string `json:"productName"`
	CurrencyID          int    `json:"currencyID"`
	CurrencyName        string `json:"currencyName"`
	CompanyID           int    `json:"companyID"`
	CompanyName         string `json:"companyName"`
	TransactionalNumber int    `json:"transactionalNumber"`
	WorkflowID          int    `json:"workflowID"`
	Maker               *User  `json:"maker"`
	UaID                int    `json:"uaID"`
}

type Records struct {
	LastUpdatedAt *timestamppb.Timestamp `json:"lastUpdatedAt"`
	TopRange      int64                  `json:"topRange"`
	BottomRange   int                    `json:"bottomRange"`
	Flows         []*Flow                `json:"flows"`
}

type Flow struct {
	WorkflowLogicID int       `json:"workflowLogicID"`
	Verifier        *Approver `json:"verifier"`
	Approver        *Approver `json:"approver"`
	Releaser        *Approver `json:"releaser"`
	VerCurrPriority int       `json:"verCurrPriority,omitempty"`
	AppCurrPriority int       `json:"appCurrPriority,omitempty"`
	ListRoles       []int     `json:"listRoles"`
}

type Approver struct {
	Requirement  int            `json:"Requirement,omitempty"`
	Participants []*Participant `json:"participants,omitempty"`
}

type Participant struct {
	RoleID         int    `json:"roleID"`
	Step           string `json:"step"`
	PriorityNumber int    `json:"priorityNumber"`
}

type User struct {
	UserID   int    `json:"userID"`
	UserName string `json:"userName,omitempty"`
	Username string `json:"username,omitempty"` // Perhatikan perbedaan field name
}
