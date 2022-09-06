package main

type Customer struct {
	ID          int    `json:"id"`
	CompanyName string `json:"companyname"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Comments    string `json:"comments"`
	Addresses   []struct {
		DefaultShipping bool   `json:"defaultshipping"`
		AddrText        string `json:"addrtext"`
	} `json:"addresses"`
	CustentityClxLocate2UId int `json:"custentity_clx_locate2uid"`
}

type Locate2UCustomer struct {
	CustomerID int `json:"customerId,omitempty"`
	TeamID     int `json:"teamId,omitempty"`

	Name      string `json:"name"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Company   string `json:"company"`
	Address   string `json:"address"`
	// Location  struct {
	// 	Latitude  string `json:"latitude"`
	// 	Longitude string `json:"longitude"`
	// } `json:"location"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Notes string `json:"notes"`
}

type Locate2UStopLine struct {
	Barcode          string `json:"barcode"`
	Description      string `json:"description"`
	CurrentLocation  string `json:"currentlocation"`
	ServiceID        int    `json:"serviceId"`
	ProductVariantID int    `json:"productVariantId"`
	Quantity         int    `json:"quantity"`
}

type Locate2UStop struct {
	StopID               int                `json:"stopId,omitempty"`
	Name                 string             `json:"name"`                 // mandatory
	Address              string             `json:"address"`              // mandatory
	Notes                string             `json:"notes"`                // mandatory
	Lines                []Locate2UStopLine `json:"lines,omitempty"`      // mandatory
	TripDate             string             `json:"tripDate"`             // mandatory
	AssignedTeamMemberID string             `json:"assignedTeamMemberId"` // mandatory
	CustomerID           int                `json:"customerId"`           // mandatory
	RunNumber            int                `json:"runNumber"`            // mandatory
	AppointmentTime      string             `json:"appointmentTime,omitempty"`
	TimeWindowStart      int                `json:"timeWindowStart,omitempty"`
	TimeWindowEnd        int                `json:"timeWindowEnd,omitempty"`
	DurationMinutes      int                `json:"durationMinutes,omitempty"`
	TeamRegionId         int                `json:"teamRegionId,omitempty"`
}

type Locate2ULink struct {
	LinkID     int    `json:"linkId,omitempty"`
	StopID     int    `json:"stopId,omitempty"`
	ShipmentID int    `json:"shipmentId,omitempty"`
	Type       string `json:"type,omitempty"`
	// ExpiryDate string `json:"expiryDate"`
	// Recipient  struct {
	// 	Name  string `json:"name"`
	// 	Phone string `json:"phone"`
	// 	Email string `json:"email"`
	// } `json:"recipient"`
	// SharingMechanism string `json:"sharingMechanism"`
	Message string `json:"message"`
	// CurrentLocation  struct {
	// 	Latitude  float32 `json:"latitude"`
	// 	Longitude float32 `json:"longitude"`
	// } `json:"currentLocation"`
	Url          string `json:"url,omitempty"`
	TeamMemberId int    `json:"teamMemberId,omitempty"`
}
