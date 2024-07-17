package us_enrichment

type PrincipalResponse struct {
	SmartyKey      string              `json:"smarty_key"`
	DataSetName    string              `json:"data_set_name"`
	DataSubsetName string              `json:"data_subset_name"`
	Attributes     PrincipalAttributes `json:"attributes"`
	Etag           string
}

type PrincipalAttributes struct {
	FirstFloorSqft                  string `json:"1st_floor_sqft"`
	SecondFloorSqft                 string `json:"2nd_floor_sqft"`
	Acres                           string `json:"acres"`
	AirConditioner                  string `json:"air_conditioner"`
	ArborPergola                    string `json:"arbor_pergola"`
	AssessedImprovementPercent      string `json:"assessed_improvement_percent"`
	AssessedImprovementValue        string `json:"assessed_improvement_value"`
	AssessedLandValue               string `json:"assessed_land_value"`
	AssessedValue                   string `json:"assessed_value"`
	AssessorLastUpdate              string `json:"assessor_last_update"`
	AssessorTaxrollUpdate           string `json:"assessor_taxroll_update"`
	AtticArea                       string `json:"attic_area"`
	AtticFlag                       string `json:"attic_flag"`
	Balcony                         string `json:"balcony"`
	BalconyArea                     string `json:"balcony_area"`
	BasementSqft                    string `json:"basement_sqft"`
	BasementSqftFinished            string `json:"basement_sqft_finished"`
	BasementSqftUnfinished          string `json:"basement_sqft_unfinished"`
	BathHouse                       string `json:"bath_house"`
	BathHouseSqft                   string `json:"bath_house_sqft"`
	BathroomsPartial                string `json:"bathrooms_partial"`
	BathroomsTotal                  string `json:"bathrooms_total"`
	Bedrooms                        string `json:"bedrooms"`
	Block1                          string `json:"block1"`
	Block2                          string `json:"block2"`
	BoatAccess                      string `json:"boat_access"`
	BoatHouse                       string `json:"boat_house"`
	BoatHouseSqft                   string `json:"boat_house_sqft"`
	BoatLift                        string `json:"boat_lift"`
	BonusRoom                       string `json:"bonus_room"`
	BreakfastNook                   string `json:"breakfast_nook"`
	Breezeway                       string `json:"breezeway"`
	BuildingDefinitionCode          string `json:"building_definition_code"`
	BuildingSqft                    string `json:"building_sqft"`
	Cabin                           string `json:"cabin"`
	CabinSqft                       string `json:"cabin_sqft"`
	Canopy                          string `json:"canopy"`
	CanopySqft                      string `json:"canopy_sqft"`
	Carport                         string `json:"carport"`
	CarportSqft                     string `json:"carport_sqft"`
	CbsaCode                        string `json:"cbsa_code"`
	CbsaName                        string `json:"cbsa_name"`
	Cellar                          string `json:"cellar"`
	CensusBlock                     string `json:"census_block"`
	CensusBlockGroup                string `json:"census_block_group"`
	CensusFipsPlaceCode             string `json:"census_fips_place_code"`
	CensusTract                     string `json:"census_tract"`
	CentralVacuum                   string `json:"central_vacuum"`
	CodeTitleCompany                string `json:"code_title_company"`
	CombinedStatisticalArea         string `json:"combined_statistical_area"`
	CommunityRec                    string `json:"community_rec"`
	CompanyFlag                     string `json:"company_flag"`
	CongressionalDistrict           string `json:"congressional_district"`
	ConstructionType                string `json:"construction_type"`
	ContactCity                     string `json:"contact_city"`
	ContactCrrt                     string `json:"contact_crrt"`
	ContactFullAddress              string `json:"contact_full_address"`
	ContactHouseNumber              string `json:"contact_house_number"`
	ContactMailInfoFormat           string `json:"contact_mail_info_format"`
	ContactMailInfoPrivacy          string `json:"contact_mail_info_privacy"`
	ContactMailingCounty            string `json:"contact_mailing_county"`
	ContactMailingFips              string `json:"contact_mailing_fips"`
	ContactPostDirection            string `json:"contact_post_direction"`
	ContactPreDirection             string `json:"contact_pre_direction"`
	ContactState                    string `json:"contact_state"`
	ContactStreetName               string `json:"contact_street_name"`
	ContactSuffix                   string `json:"contact_suffix"`
	ContactUnitDesignator           string `json:"contact_unit_designator"`
	ContactValue                    string `json:"contact_value"`
	ContactZip                      string `json:"contact_zip"`
	ContactZip4                     string `json:"contact_zip4"`
	Courtyard                       string `json:"courtyard"`
	CourtyardArea                   string `json:"courtyard_area"`
	Deck                            string `json:"deck"`
	DeckArea                        string `json:"deck_area"`
	DeedDocumentPage                string `json:"deed_document_page"`
	DeedDocumentBook                string `json:"deed_document_book"`
	DeedDocumentNumber              string `json:"deed_document_number"`
	DeedOwnerFirstName              string `json:"deed_owner_first_name"`
	DeedOwnerFirstName2             string `json:"deed_owner_first_name2"`
	DeedOwnerFirstName3             string `json:"deed_owner_first_name3"`
	DeedOwnerFirstName4             string `json:"deed_owner_first_name4"`
	DeedOwnerFullName               string `json:"deed_owner_full_name"`
	DeedOwnerFullName2              string `json:"deed_owner_full_name2"`
	DeedOwnerFullName3              string `json:"deed_owner_full_name3"`
	DeedOwnerFullName4              string `json:"deed_owner_full_name4"`
	DeedOwnerLastName               string `json:"deed_owner_last_name"`
	DeedOwnerLastName2              string `json:"deed_owner_last_name2"`
	DeedOwnerLastName3              string `json:"deed_owner_last_name3"`
	DeedOwnerLastName4              string `json:"deed_owner_last_name4"`
	DeedOwnerMiddleName             string `json:"deed_owner_middle_name"`
	DeedOwnerMiddleName2            string `json:"deed_owner_middle_name2"`
	DeedOwnerMiddleName3            string `json:"deed_owner_middle_name3"`
	DeedOwnerMiddleName4            string `json:"deed_owner_middle_name4"`
	DeedOwnerSuffix                 string `json:"deed_owner_suffix"`
	DeedOwnerSuffix2                string `json:"deed_owner_suffix2"`
	DeedOwnerSuffix3                string `json:"deed_owner_suffix3"`
	DeedOwnerSuffix4                string `json:"deed_owner_suffix4"`
	DeedSaleDate                    string `json:"deed_sale_date"`
	DeedSalePrice                   string `json:"deed_sale_price"`
	DeedTransactionId               string `json:"deed_transaction_id"`
	DepthLinearFootage              string `json:"depth_linear_footage"`
	DisabledTaxExemption            string `json:"disabled_tax_exemption"`
	DocumentTypeDescription         string `json:"document_type_description"`
	DrivewaySqft                    string `json:"driveway_sqft"`
	DrivewayType                    string `json:"driveway_type"`
	EffectiveYearBuilt              string `json:"effective_year_built"`
	ElevationFeet                   string `json:"elevation_feet"`
	Elevator                        string `json:"elevator"`
	EquestrianArena                 string `json:"equestrian_arena"`
	Escalator                       string `json:"escalator"`
	ExerciseRoom                    string `json:"exercise_room"`
	ExteriorWalls                   string `json:"exterior_walls"`
	FamilyRoom                      string `json:"family_room"`
	Fence                           string `json:"fence"`
	FenceArea                       string `json:"fence_area"`
	FipsCode                        string `json:"fips_code"`
	FireResistanceCode              string `json:"fire_resistance_code"`
	FireSprinklersFlag              string `json:"fire_sprinklers_flag"`
	Fireplace                       string `json:"fireplace"`
	FireplaceNumber                 string `json:"fireplace_number"`
	FirstName                       string `json:"first_name"`
	FirstName2                      string `json:"first_name_2"`
	FirstName3                      string `json:"first_name_3"`
	FirstName4                      string `json:"first_name_4"`
	Flooring                        string `json:"flooring"`
	Foundation                      string `json:"foundation"`
	GameRoom                        string `json:"game_room"`
	Garage                          string `json:"garage"`
	GarageSqft                      string `json:"garage_sqft"`
	Gazebo                          string `json:"gazebo"`
	GazeboSqft                      string `json:"gazebo_sqft"`
	GolfCourse                      string `json:"golf_course"`
	Grainery                        string `json:"grainery"`
	GrainerySqft                    string `json:"grainery_sqft"`
	GreatRoom                       string `json:"great_room"`
	Greenhouse                      string `json:"greenhouse"`
	GreenhouseSqft                  string `json:"greenhouse_sqft"`
	GrossSqft                       string `json:"gross_sqft"`
	Guesthouse                      string `json:"guesthouse"`
	GuesthouseSqft                  string `json:"guesthouse_sqft"`
	HandicapAccessibility           string `json:"handicap_accessibility"`
	Heat                            string `json:"heat"`
	HeatFuelType                    string `json:"heat_fuel_type"`
	HobbyRoom                       string `json:"hobby_room"`
	HomeownerTaxExemption           string `json:"homeowner_tax_exemption"`
	InstrumentDate                  string `json:"instrument_date"`
	IntercomSystem                  string `json:"intercom_system"`
	InterestRateType2               string `json:"interest_rate_type_2"`
	InteriorStructure               string `json:"interior_structure"`
	Kennel                          string `json:"kennel"`
	KennelSqft                      string `json:"kennel_sqft"`
	LandUseCode                     string `json:"land_use_code"`
	LandUseGroup                    string `json:"land_use_group"`
	LandUseStandard                 string `json:"land_use_standard"`
	LastName                        string `json:"last_name"`
	LastName2                       string `json:"last_name_2"`
	LastName3                       string `json:"last_name_3"`
	LastName4                       string `json:"last_name_4"`
	Latitude                        string `json:"latitude"`
	Laundry                         string `json:"laundry"`
	LeanTo                          string `json:"lean_to"`
	LeanToSqft                      string `json:"lean_to_sqft"`
	LegalDescription                string `json:"legal_description"`
	LegalUnit                       string `json:"legal_unit"`
	LenderAddress                   string `json:"lender_address"`
	LenderAddress2                  string `json:"lender_address_2"`
	LenderCity                      string `json:"lender_city"`
	LenderCity2                     string `json:"lender_city_2"`
	LenderCode2                     string `json:"lender_code_2"`
	LenderFirstName                 string `json:"lender_first_name"`
	LenderFirstName2                string `json:"lender_first_name_2"`
	LenderLastName                  string `json:"lender_last_name"`
	LenderLastName2                 string `json:"lender_last_name_2"`
	LenderName                      string `json:"lender_name"`
	LenderName2                     string `json:"lender_name_2"`
	LenderSellerCarryBack           string `json:"lender_seller_carry_back"`
	LenderSellerCarryBack2          string `json:"lender_seller_carry_back_2"`
	LenderState                     string `json:"lender_state"`
	LenderState2                    string `json:"lender_state_2"`
	LenderZip                       string `json:"lender_zip"`
	LenderZip2                      string `json:"lender_zip_2"`
	LenderZipExtended               string `json:"lender_zip_extended"`
	LenderZipExtended2              string `json:"lender_zip_extended_2"`
	LoadingPlatform                 string `json:"loading_platform"`
	LoadingPlatformSqft             string `json:"loading_platform_sqft"`
	Longitude                       string `json:"longitude"`
	Lot1                            string `json:"lot_1"`
	Lot2                            string `json:"lot_2"`
	Lot3                            string `json:"lot_3"`
	LotSqft                         string `json:"lot_sqft"`
	MarketImprovementPercent        string `json:"market_improvement_percent"`
	MarketImprovementValue          string `json:"market_improvement_value"`
	MarketLandValue                 string `json:"market_land_value"`
	MarketValueYear                 string `json:"market_value_year"`
	MatchType                       string `json:"match_type"`
	MediaRoom                       string `json:"media_room"`
	MetroDivision                   string `json:"metro_division"`
	MiddleName                      string `json:"middle_name"`
	MiddleName2                     string `json:"middle_name_2"`
	MiddleName3                     string `json:"middle_name_3"`
	MiddleName4                     string `json:"middle_name_4"`
	Milkhouse                       string `json:"milkhouse"`
	MilkhouseSqft                   string `json:"milkhouse_sqft"`
	MinorCivilDivisionCode          string `json:"minor_civil_division_code"`
	MinorCivilDivisionName          string `json:"minor_civil_division_name"`
	MobileHomeHookup                string `json:"mobile_home_hookup"`
	MortgageAmount                  string `json:"mortgage_amount"`
	MortgageAmount2                 string `json:"mortgage_amount_2"`
	MortgageDueDate                 string `json:"mortgage_due_date"`
	MortgageDueDate2                string `json:"mortgage_due_date_2"`
	MortgageInterestRate            string `json:"mortgage_interest_rate"`
	MortgageInterestRateType        string `json:"mortgage_interest_rate_type"`
	MortgageLenderCode              string `json:"mortgage_lender_code"`
	MortgageRate2                   string `json:"mortgage_rate_2"`
	MortgageRecordingDate           string `json:"mortgage_recording_date"`
	MortgageRecordingDate2          string `json:"mortgage_recording_date_2"`
	MortgageTerm                    string `json:"mortgage_term"`
	MortgageTerm2                   string `json:"mortgage_term_2"`
	MortgageTermType                string `json:"mortgage_term_type"`
	MortgageTermType2               string `json:"mortgage_term_type_2"`
	MortgageType                    string `json:"mortgage_type"`
	MortgageType2                   string `json:"mortgage_type_2"`
	MsaCode                         string `json:"msa_code"`
	MsaName                         string `json:"msa_name"`
	MudRoom                         string `json:"mud_room"`
	MultiParcelFlag                 string `json:"multi_parcel_flag"`
	NameTitleCompany                string `json:"name_title_company"`
	NeighborhoodCode                string `json:"neighborhood_code"`
	NumberOfBuildings               string `json:"number_of_buildings"`
	Office                          string `json:"office"`
	OfficeSqft                      string `json:"office_sqft"`
	OtherTaxExemption               string `json:"other_tax_exemption"`
	OutdoorKitchenFireplace         string `json:"outdoor_kitchen_fireplace"`
	OverheadDoor                    string `json:"overhead_door"`
	OwnerFullName                   string `json:"owner_full_name"`
	OwnerFullName2                  string `json:"owner_full_name_2"`
	OwnerFullName3                  string `json:"owner_full_name_3"`
	OwnerFullName4                  string `json:"owner_full_name_4"`
	OwnerOccupancyStatus            string `json:"owner_occupancy_status"`
	OwnershipTransferDate           string `json:"ownership_transfer_date"`
	OwnershipTransferDocNumber      string `json:"ownership_transfer_doc_number"`
	OwnershipTransferTransactionId  string `json:"ownership_transfer_transaction_id"`
	OwnershipType                   string `json:"ownership_type"`
	OwnershipType2                  string `json:"ownership_type_2"`
	OwnershipVestingRelationCode    string `json:"ownership_vesting_relation_code"`
	ParcelAccountNumber             string `json:"parcel_account_number"`
	ParcelMapBook                   string `json:"parcel_map_book"`
	ParcelMapPage                   string `json:"parcel_map_page"`
	ParcelNumberAlternate           string `json:"parcel_number_alternate"`
	ParcelNumberFormatted           string `json:"parcel_number_formatted"`
	ParcelNumberPrevious            string `json:"parcel_number_previous"`
	ParcelNumberYearAdded           string `json:"parcel_number_year_added"`
	ParcelNumberYearChange          string `json:"parcel_number_year_change"`
	ParcelRawNumber                 string `json:"parcel_raw_number"`
	ParcelShellRecord               string `json:"parcel_shell_record"`
	ParkingSpaces                   string `json:"parking_spaces"`
	PatioArea                       string `json:"patio_area"`
	PhaseName                       string `json:"phase_name"`
	PlumbingFixturesCount           string `json:"plumbing_fixtures_count"`
	PoleStruct                      string `json:"pole_struct"`
	PoleStructSqft                  string `json:"pole_struct_sqft"`
	Pond                            string `json:"pond"`
	Pool                            string `json:"pool"`
	PoolArea                        string `json:"pool_area"`
	Poolhouse                       string `json:"poolhouse"`
	PoolhouseSqft                   string `json:"poolhouse_sqft"`
	Porch                           string `json:"porch"`
	PorchArea                       string `json:"porch_area"`
	PoultryHouse                    string `json:"poultry_house"`
	PoultryHouseSqft                string `json:"poultry_house_sqft"`
	PreviousAssessedValue           string `json:"previous_assessed_value"`
	PriorSaleAmount                 string `json:"prior_sale_amount"`
	PriorSaleDate                   string `json:"prior_sale_date"`
	PropertyAddressCarrierRouteCode string `json:"property_address_carrier_route_code"`
	PropertyAddressCity             string `json:"property_address_city"`
	PropertyAddressFull             string `json:"property_address_full"`
	PropertyAddressHouseNumber      string `json:"property_address_house_number"`
	PropertyAddressPostDirection    string `json:"property_address_post_direction"`
	PropertyAddressPreDirection     string `json:"property_address_pre_direction"`
	PropertyAddressState            string `json:"property_address_state"`
	PropertyAddressStreetName       string `json:"property_address_street_name"`
	PropertyAddressStreetSuffix     string `json:"property_address_street_suffix"`
	PropertyAddressUnitDesignator   string `json:"property_address_unit_designator"`
	PropertyAddressUnitValue        string `json:"property_address_unit_value"`
	PropertyAddressZip4             string `json:"property_address_zip_4"`
	PropertyAddressZipcode          string `json:"property_address_zipcode"`
	PublicationDate                 string `json:"publication_date"`
	Quarter                         string `json:"quarter"`
	QuarterQuarter                  string `json:"quarter_quarter"`
	Quonset                         string `json:"quonset"`
	QuonsetSqft                     string `json:"quonset_sqft"`
	Range                           string `json:"range"`
	RecordingDate                   string `json:"recording_date"`
	RoofCover                       string `json:"roof_cover"`
	RoofFrame                       string `json:"roof_frame"`
	Rooms                           string `json:"rooms"`
	RvParking                       string `json:"rv_parking"`
	SafeRoom                        string `json:"safe_room"`
	SaleAmount                      string `json:"sale_amount"`
	SaleDate                        string `json:"sale_date"`
	Sauna                           string `json:"sauna"`
	Section                         string `json:"section"`
	SecurityAlarm                   string `json:"security_alarm"`
	SeniorTaxExemption              string `json:"senior_tax_exemption"`
	SewerType                       string `json:"sewer_type"`
	Shed                            string `json:"shed"`
	ShedSqft                        string `json:"shed_sqft"`
	Silo                            string `json:"silo"`
	SiloSqft                        string `json:"silo_sqft"`
	SittingRoom                     string `json:"sitting_room"`
	SitusCounty                     string `json:"situs_county"`
	SitusState                      string `json:"situs_state"`
	SoundSystem                     string `json:"sound_system"`
	SportsCourt                     string `json:"sports_court"`
	Sprinklers                      string `json:"sprinklers"`
	Stable                          string `json:"stable"`
	StableSqft                      string `json:"stable_sqft"`
	StorageBuilding                 string `json:"storage_building"`
	StorageBuildingSqft             string `json:"storage_building_sqft"`
	StoriesNumber                   string `json:"stories_number"`
	StormShelter                    string `json:"storm_shelter"`
	StormShutter                    string `json:"storm_shutter"`
	StructureStyle                  string `json:"structure_style"`
	Study                           string `json:"study"`
	Subdivision                     string `json:"subdivision"`
	Suffix                          string `json:"suffix"`
	Suffix2                         string `json:"suffix_2"`
	Suffix3                         string `json:"suffix_3"`
	Suffix4                         string `json:"suffix_4"`
	Sunroom                         string `json:"sunroom"`
	TaxAssessYear                   string `json:"tax_assess_year"`
	TaxBilledAmount                 string `json:"tax_billed_amount"`
	TaxDelinquentYear               string `json:"tax_delinquent_year"`
	TaxFiscalYear                   string `json:"tax_fiscal_year"`
	TaxJurisdiction                 string `json:"tax_jurisdiction"`
	TaxRateArea                     string `json:"tax_rate_area"`
	TennisCourt                     string `json:"tennis_court"`
	TopographyCode                  string `json:"topography_code"`
	TotalMarketValue                string `json:"total_market_value"`
	Township                        string `json:"township"`
	TractNumber                     string `json:"tract_number"`
	TransferAmount                  string `json:"transfer_amount"`
	TrustDescription                string `json:"trust_description"`
	UnitCount                       string `json:"unit_count"`
	UpperFloorsSqft                 string `json:"upper_floors_sqft"`
	Utility                         string `json:"utility"`
	UtilityBuilding                 string `json:"utility_building"`
	UtilityBuildingSqft             string `json:"utility_building_sqft"`
	UtilitySqft                     string `json:"utility_sqft"`
	VeteranTaxExemption             string `json:"veteran_tax_exemption"`
	ViewDescription                 string `json:"view_description"`
	WaterFeature                    string `json:"water_feature"`
	WaterServiceType                string `json:"water_service_type"`
	WetBar                          string `json:"wet_bar"`
	WidowTaxExemption               string `json:"widow_tax_exemption"`
	WidthLinearFootage              string `json:"width_linear_footage"`
	WineCellar                      string `json:"wine_cellar"`
	YearBuilt                       string `json:"year_built"`
	Zoning                          string `json:"zoning"`
}

type FinancialResponse struct {
	SmartyKey      string              `json:"smarty_key"`
	DataSetName    string              `json:"data_set_name"`
	DataSubsetName string              `json:"data_subset_name"`
	Attributes     FinancialAttributes `json:"attributes"`
	Etag           string
}

type FinancialAttributes struct {
	AssessedImprovementPercent string `json:"assessed_improvement_percent"`
	AssessedImprovementValue   string `json:"assessed_improvement_value"`
	AssessedLandValue          string `json:"assessed_land_value"`
	AssessedValue              string `json:"assessed_value"`
	AssessorLastUpdate         string `json:"assessor_last_update"`
	AssessorTaxrollUpdate      string `json:"assessor_taxroll_update"`
	ContactCity                string `json:"contact_city"`
	ContactCrrt                string `json:"contact_crrt"`
	ContactFullAddress         string `json:"contact_full_address"`
	ContactHouseNumber         string `json:"contact_house_number"`
	ContactMailInfoFormat      string `json:"contact_mail_info_format"`
	ContactMailInfoPrivacy     string `json:"contact_mail_info_privacy"`
	ContactMailingCounty       string `json:"contact_mailing_county"`
	ContactMailingFips         string `json:"contact_mailing_fips"`
	ContactPostDirection       string `json:"contact_post_direction"`
	ContactPreDirection        string `json:"contact_pre_direction"`
	ContactState               string `json:"contact_state"`
	ContactStreetName          string `json:"contact_street_name"`
	ContactSuffix              string `json:"contact_suffix"`
	ContactUnitDesignator      string `json:"contact_unit_designator"`
	ContactValue               string `json:"contact_value"`
	ContactZip                 string `json:"contact_zip"`
	ContactZip4                string `json:"contact_zip4"`
	DeedDocumentPage           string `json:"deed_document_page"`
	DeedDocumentBook           string `json:"deed_document_book"`
	DeedDocumentNumber         string `json:"deed_document_number"`
	DeedOwnerFirstName         string `json:"deed_owner_first_name"`
	DeedOwnerFirstName2        string `json:"deed_owner_first_name2"`
	DeedOwnerFirstName3        string `json:"deed_owner_first_name3"`
	DeedOwnerFirstName4        string `json:"deed_owner_first_name4"`
	DeedOwnerFullName          string `json:"deed_owner_full_name"`
	DeedOwnerFullName2         string `json:"deed_owner_full_name2"`
	DeedOwnerFullName3         string `json:"deed_owner_full_name3"`
	DeedOwnerFullName4         string `json:"deed_owner_full_name4"`
	DeedOwnerLastName          string `json:"deed_owner_last_name"`
	DeedOwnerLastName2         string `json:"deed_owner_last_name2"`
	DeedOwnerLastName3         string `json:"deed_owner_last_name3"`
	DeedOwnerLastName4         string `json:"deed_owner_last_name4"`
	DeedOwnerMiddleName        string `json:"deed_owner_middle_name"`
	DeedOwnerMiddleName2       string `json:"deed_owner_middle_name2"`
	DeedOwnerMiddleName3       string `json:"deed_owner_middle_name3"`
	DeedOwnerMiddleName4       string `json:"deed_owner_middle_name4"`
	DeedOwnerSuffix            string `json:"deed_owner_suffix"`
	DeedOwnerSuffix2           string `json:"deed_owner_suffix2"`
	DeedOwnerSuffix3           string `json:"deed_owner_suffix3"`
	DeedOwnerSuffix4           string `json:"deed_owner_suffix4"`
	DeedSaleDate               string `json:"deed_sale_date"`
	DeedSalePrice              string `json:"deed_sale_price"`
	DeedTransactionId          string `json:"deed_transaction_id"`
	DisabledTaxExemption       string `json:"disabled_tax_exemption"`
	FinancialHistory           []struct {
		CodeTitleCompany         string `json:"code_title_company"`
		DocumentTypeDescription  string `json:"document_type_description"`
		InstrumentDate           string `json:"instrument_date"`
		InterestRateType2        string `json:"interest_rate_type_2"`
		LenderAddress            string `json:"lender_address"`
		LenderAddress2           string `json:"lender_address_2"`
		LenderCity               string `json:"lender_city"`
		LenderCity2              string `json:"lender_city_2"`
		LenderCode2              string `json:"lender_code_2"`
		LenderFirstName          string `json:"lender_first_name"`
		LenderFirstName2         string `json:"lender_first_name_2"`
		LenderLastName           string `json:"lender_last_name"`
		LenderLastName2          string `json:"lender_last_name_2"`
		LenderName               string `json:"lender_name"`
		LenderName2              string `json:"lender_name_2"`
		LenderSellerCarryBack    string `json:"lender_seller_carry_back"`
		LenderSellerCarryBack2   string `json:"lender_seller_carry_back_2"`
		LenderState              string `json:"lender_state"`
		LenderState2             string `json:"lender_state_2"`
		LenderZip                string `json:"lender_zip"`
		LenderZip2               string `json:"lender_zip_2"`
		LenderZipExtended        string `json:"lender_zip_extended"`
		LenderZipExtended2       string `json:"lender_zip_extended_2"`
		MortgageAmount           string `json:"mortgage_amount"`
		MortgageAmount2          string `json:"mortgage_amount_2"`
		MortgageDueDate          string `json:"mortgage_due_date"`
		MortgageDueDate2         string `json:"mortgage_due_date_2"`
		MortgageInterestRate     string `json:"mortgage_interest_rate"`
		MortgageInterestRateType string `json:"mortgage_interest_rate_type"`
		MortgageLenderCode       string `json:"mortgage_lender_code"`
		MortgageRate2            string `json:"mortgage_rate_2"`
		MortgageRecordingDate    string `json:"mortgage_recording_date"`
		MortgageRecordingDate2   string `json:"mortgage_recording_date_2"`
		MortgageTerm             string `json:"mortgage_term"`
		MortgageTerm2            string `json:"mortgage_term_2"`
		MortgageTermType         string `json:"mortgage_term_type"`
		MortgageTermType2        string `json:"mortgage_term_type_2"`
		MortgageType             string `json:"mortgage_type"`
		MortgageType2            string `json:"mortgage_type_2"`
		MultiParcelFlag          string `json:"multi_parcel_flag"`
		NameTitleCompany         string `json:"name_title_company"`
		RecordingDate            string `json:"recording_date"`
		TransferAmount           string `json:"transfer_amount"`
	} `json:"financial_history"`
	FirstName                      string `json:"first_name"`
	FirstName2                     string `json:"first_name_2"`
	FirstName3                     string `json:"first_name_3"`
	FirstName4                     string `json:"first_name_4"`
	HomeownerTaxExemption          string `json:"homeowner_tax_exemption"`
	LastName                       string `json:"last_name"`
	LastName2                      string `json:"last_name_2"`
	LastName3                      string `json:"last_name_3"`
	LastName4                      string `json:"last_name_4"`
	MarketImprovementPercent       string `json:"market_improvement_percent"`
	MarketImprovementValue         string `json:"market_improvement_value"`
	MarketLandValue                string `json:"market_land_value"`
	MarketValueYear                string `json:"market_value_year"`
	MatchType                      string `json:"match_type"`
	MiddleName                     string `json:"middle_name"`
	MiddleName2                    string `json:"middle_name_2"`
	MiddleName3                    string `json:"middle_name_3"`
	MiddleName4                    string `json:"middle_name_4"`
	OtherTaxExemption              string `json:"other_tax_exemption"`
	OwnerFullName                  string `json:"owner_full_name"`
	OwnerFullName2                 string `json:"owner_full_name_2"`
	OwnerFullName3                 string `json:"owner_full_name_3"`
	OwnerFullName4                 string `json:"owner_full_name_4"`
	OwnershipTransferDate          string `json:"ownership_transfer_date"`
	OwnershipTransferDocNumber     string `json:"ownership_transfer_doc_number"`
	OwnershipTransferTransactionId string `json:"ownership_transfer_transaction_id"`
	OwnershipType                  string `json:"ownership_type"`
	OwnershipType2                 string `json:"ownership_type_2"`
	PreviousAssessedValue          string `json:"previous_assessed_value"`
	PriorSaleAmount                string `json:"prior_sale_amount"`
	PriorSaleDate                  string `json:"prior_sale_date"`
	SaleAmount                     string `json:"sale_amount"`
	SaleDate                       string `json:"sale_date"`
	SeniorTaxExemption             string `json:"senior_tax_exemption"`
	Suffix                         string `json:"suffix"`
	Suffix2                        string `json:"suffix_2"`
	Suffix3                        string `json:"suffix_3"`
	Suffix4                        string `json:"suffix_4"`
	TaxAssessYear                  string `json:"tax_assess_year"`
	TaxBilledAmount                string `json:"tax_billed_amount"`
	TaxDelinquentYear              string `json:"tax_delinquent_year"`
	TaxFiscalYear                  string `json:"tax_fiscal_year"`
	TaxRateArea                    string `json:"tax_rate_area"`
	TotalMarketValue               string `json:"total_market_value"`
	TrustDescription               string `json:"trust_description"`
	VeteranTaxExemption            string `json:"veteran_tax_exemption"`
	WidowTaxExemption              string `json:"widow_tax_exemption"`
}

type GeoReferenceResponse struct {
	SmartyKey   string                 `json:"smarty_key"`
	DataSetName string                 `json:"data_set_name"`
	Attributes  GeoReferenceAttributes `json:"attributes"`
	Etag        string
}

type GeoReferenceAttributes struct {
	CensusBlock struct {
		Accuracy string `json:"accuracy"`
		Id       string `json:"id"`
	} `json:"census_block"`

	CensusCountyDivision struct {
		Accuracy string `json:"accuracy"`
		Id       string `json:"id"`
		Name     string `json:"name"`
	} `json:"census_county_division"`

	CoreBasedStatArea struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"core_based_stat_area"`

	Place struct {
		Accuracy string `json:"accuracy"`
		Id       string `json:"id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
	} `json:"place"`
}

type RootAddress struct {
	SecondaryCount      uint64 `json:"secondary_count"`
	SmartyKey           string `json:"smarty_key"`
	PrimaryNumber       string `json:"primary_number"`
	StreetPredirection  string `json:"street_predirection"`
	StreetName          string `json:"street_name"`
	StreetSuffix        string `json:"street_suffix"`
	StreetPostdirection string `json:"street_postdirection"`
	CityName            string `json:"city_name"`
	StateAbbreviation   string `json:"state_abbreviation"`
	Zipcode             string `json:"zipcode"`
	Plus4Code           string `json:"plus4_code"`
}

type Alias struct {
	SmartyKey           string `json:"smarty_key"`
	PrimaryNumber       string `json:"primary_number"`
	StreetPredirection  string `json:"street_predirection"`
	StreetName          string `json:"street_name"`
	StreetSuffix        string `json:"street_suffix"`
	StreetPostdirection string `json:"street_postdirection"`
	CityName            string `json:"city_name"`
	StateAbbreviation   string `json:"state_abbreviation"`
	Zipcode             string `json:"zipcode"`
	Plus4Code           string `json:"plus4_code"`
}

type Secondary struct {
	SmartyKey           string `json:"smarty_key"`
	SecondaryDesignator string `json:"secondary_designator"`
	SecondaryNumber     string `json:"secondary_number"`
	Plus4Code           string `json:"plus4_code"`
}

type SecondaryResponse struct {
	SmartyKey   string `json:"smarty_key"`
	Etag        string
	RootAddress RootAddress `json:"root_address"`
	Aliases     []Alias     `json:"aliases"`
	Secondaries []Secondary `json:"secondaries"`
}

type SecondaryCountResponse struct {
	SmartyKey string `json:"smarty_key"`
	Count     int    `json:"count"`
	Etag      string
}
