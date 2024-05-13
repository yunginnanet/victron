package victron

type (
	ProductID int
	Category  int
)

const (
	CategoryUnknown Category = iota
	CategoryBattery
	CategorySolar
	CategoryInverter
	CategoryController
	CategoryOther
)

// BEGIN PRODUCT IDS
const (
	PIDBMV700                       ProductID = 0x203
	PIDBMV702                       ProductID = 0x204
	PIDBMV700H                      ProductID = 0x205
	PIDBlueSolarMPPT7015            ProductID = 0x300
	PIDBlueSolarMPPT7550            ProductID = 0xA040
	PIDBlueSolarMPPT15035           ProductID = 0xA041
	PIDBlueSolarMPPT7515            ProductID = 0xA042
	PIDBlueSolarMPPT10015           ProductID = 0xA043
	PIDBlueSolarMPPT10030           ProductID = 0xA044
	PIDBlueSolarMPPT10050           ProductID = 0xA045
	PIDBlueSolarMPPT15070           ProductID = 0xA046
	PIDBlueSolarMPPT150100          ProductID = 0xA047
	PIDSmartSolarMPPT250100         ProductID = 0xA050
	PIDSmartSolarMPPT150100         ProductID = 0xA051
	PIDSmartSolarMPPT15085          ProductID = 0xA052
	PIDSmartSolarMPPT7515           ProductID = 0xA053
	PIDSmartSolarMPPT7510           ProductID = 0xA054
	PIDSmartSolarMPPT10015          ProductID = 0xA055
	PIDSmartSolarMPPT10030          ProductID = 0xA056
	PIDSmartSolarMPPT10050          ProductID = 0xA057
	PIDSmartSolarMPPT15035          ProductID = 0xA058
	PIDSmartSolarMPPT150100Rev2     ProductID = 0xA059
	PIDSmartSolarMPPT15085Rev2      ProductID = 0xA05A
	PIDSmartSolarMPPT25070          ProductID = 0xA05B
	PIDSmartSolarMPPT25085          ProductID = 0xA05C
	PIDSmartSolarMPPT25060          ProductID = 0xA05D
	PIDSmartSolarMPPT25045          ProductID = 0xA05E
	PIDSmartSolarMPPT10020          ProductID = 0xA05F
	PIDSmartSolarMPPT1002048V       ProductID = 0xA060
	PIDSmartSolarMPPT15045          ProductID = 0xA061
	PIDSmartSolarMPPT15060          ProductID = 0xA062
	PIDSmartSolarMPPT15070          ProductID = 0xA063
	PIDSmartSolarMPPT25085Rev2      ProductID = 0xA064
	PIDSmartSolarMPPT250100Rev2     ProductID = 0xA065
	PIDBlueSolarMPPT10020           ProductID = 0xA066
	PIDBlueSolarMPPT1002048V        ProductID = 0xA067
	PIDSmartSolarMPPT25060Rev2      ProductID = 0xA068
	PIDSmartSolarMPPT25070Rev2      ProductID = 0xA069
	PIDSmartSolarMPPT15045Rev2      ProductID = 0xA06A
	PIDSmartSolarMPPT15060Rev2      ProductID = 0xA06B
	PIDSmartSolarMPPT15070Rev2      ProductID = 0xA06C
	PIDSmartSolarMPPT15085Rev3      ProductID = 0xA06D
	PIDSmartSolarMPPT150100Rev3     ProductID = 0xA06E
	PIDBlueSolarMPPT15045Rev2       ProductID = 0xA06F
	PIDBlueSolarMPPT15060Rev2       ProductID = 0xA070
	PIDBlueSolarMPPT15070Rev2       ProductID = 0xA071
	PIDPhoenixInverter12V250VA230V  ProductID = 0xA201
	PIDPhoenixInverter24V250VA230V  ProductID = 0xA202
	PIDPhoenixInverter48V250VA230V  ProductID = 0xA204
	PIDPhoenixInverter12V375VA230V  ProductID = 0xA211
	PIDPhoenixInverter24V375VA230V  ProductID = 0xA212
	PIDPhoenixInverter48V375VA230V  ProductID = 0xA214
	PIDPhoenixInverter12V500VA230V  ProductID = 0xA221
	PIDPhoenixInverter24V500VA230V  ProductID = 0xA222
	PIDPhoenixInverter48V500VA230V  ProductID = 0xA224
	PIDPhoenixInverter12V250VA120V  ProductID = 0xA239
	PIDPhoenixInverter24V250VA120V  ProductID = 0xA23A
	PIDPhoenixInverter48V250VA120V  ProductID = 0xA23C
	PIDPhoenixInverter12V375VA120V  ProductID = 0xA249
	PIDPhoenixInverter24V375VA120V  ProductID = 0xA24A
	PIDPhoenixInverter48V375VA120V  ProductID = 0xA24C
	PIDPhoenixInverter12V500VA120V  ProductID = 0xA259
	PIDPhoenixInverter24V500VA120V  ProductID = 0xA25A
	PIDPhoenixInverter48V500VA120V  ProductID = 0xA25C
	PIDPhoenixInverter12V800VA230V  ProductID = 0xA261
	PIDPhoenixInverter24V800VA230V  ProductID = 0xA262
	PIDPhoenixInverter48V800VA230V  ProductID = 0xA264
	PIDPhoenixInverter12V800VA120V  ProductID = 0xA269
	PIDPhoenixInverter24V800VA120V  ProductID = 0xA26A
	PIDPhoenixInverter48V800VA120V  ProductID = 0xA26C
	PIDPhoenixInverter12V1200VA230V ProductID = 0xA271
	PIDPhoenixInverter24V1200VA230V ProductID = 0xA272
	PIDPhoenixInverter48V1200VA230V ProductID = 0xA274
	PIDPhoenixInverter12V1200VA120V ProductID = 0xA279
	PIDPhoenixInverter24V1200VA120V ProductID = 0xA27A
	PIDPhoenixInverter48V1200VA120V ProductID = 0xA27C
	PIDPhoenixInverter12V1600VA230V ProductID = 0xA281
	PIDPhoenixInverter24V1600VA230V ProductID = 0xA282
	PIDPhoenixInverter48V1600VA230V ProductID = 0xA284
	PIDPhoenixInverter12V2000VA230V ProductID = 0xA291
	PIDPhoenixInverter24V2000VA230V ProductID = 0xA292
	PIDPhoenixInverter48V2000VA230V ProductID = 0xA294
	PIDPhoenixInverter12V3000VA230V ProductID = 0xA2A1
	PIDPhoenixInverter24V3000VA230V ProductID = 0xA2A2
	PIDPhoenixInverter48V3000VA230V ProductID = 0xA2A4
	PIDPhoenixSmartIP431250         ProductID = 0xA340
	PIDPhoenixSmartIP433250         ProductID = 0xA341
	PIDPhoenixSmartIP431230         ProductID = 0xA344
	PIDPhoenixSmartIP433230         ProductID = 0xA345
	PIDPhoenixSmartIP431616         ProductID = 0xA346
	PIDPhoenixSmartIP433616         ProductID = 0xA347
	PIDBMV712Smart                  ProductID = 0xA381
	PIDBMV710HSmart                 ProductID = 0xA382
	PIDBMV712SmartRev2              ProductID = 0xA383
	PIDSmartShunt500A50mV           ProductID = 0xA389
	PIDSmartShunt1000A50mV          ProductID = 0xA38A
	PIDSmartShunt2000A50mV          ProductID = 0xA38B
	PIDSmartBuckBoost12V12V50A      ProductID = 0xA3F0
)

// END PRODUCT IDS
