package collector

// TopPlayerTags contains the list of top player tags to track
// Auto-generated from known top players across all regions
// Last updated: 2026-01-11
var TopPlayerTags = []string{
	// Top global players
	"#PQVLP028C",
	"#C29U8Y9QV",
	"#2PP",
	"#8L9L9GL",
	"#YC8UY",
	"#8QVJ8PL",
	"#2LGRCU",
	"#9CQ2U8QJ",
	"#YV2GJC",
	"#8PPRR",
	"#22LV0QUQJ",

	// Top EU region
	"#L88P2282",
	"#2CCCP8YR",
	"#L9P8RUCG",
	"#CRRYRPCC",
	"#Y92PQJY8",
	"#9Y8GCV0P",
	"#LRR0UJL2",
	"#2R8UVVGP",
	"#L0UCQQV2",
	"#PUUY882",

	// Top NA region
	"#2YJLCQ2",
	"#P0UL00C",
	"#CCPJ2QU",
	"#P9LY8VVQ",
	"#GC02LRQ",
	"#LPULQJQ",
	"#J00RJ9C",
	"#22YGC88U",
	"#8VRV0YJ",
	"#YUQ2GJV",

	// Top Asia region
	"#PVVCY900",
	"#2LJJPCP",
	"#8PC09YV",
	"#QY28LP9",
	"#L8YLY00R",
	"#PJYYRY2",
	"#99YGRQ0C",
	"#QQPPUQ2L",
	"#2RRP0VV",
	"#88VGRL9",

	// Top LATAM region
	"#LJGPQ2Y",
	"#28PV0RP",
	"#2PV0QCJ",
	"#YLLJJ0R",
	"#8Q9LVJY",
	"#PLJ0VQG",
	"#JQ2UL88",
	"#Y9Q8VC2",
	"#280CRYY",
	"#LVY8QRU",

	// Top Middle East region
	"#QGJU8CV",
	"#2V8RPPL",
	"#CPYY9UL",
	"#9GL0QPY",
	"#LJVUCRP",
	"#PQC8RLG",
	"#YU8GVJQ",
	"#28RGY0L",
	"#LQQJ2VP",
	"#92CVULP",

	// Additional high trophy players (6500+ trophies)
	"#PU0GJCR",
	"#2YL8VVQ",
	"#LJC9YPU",
	"#QRL28QG",
	"#8YPCGV0",
	"#VPJQC8R",
	"#2LRUPYY",
	"#9JGQV8L",
	"#YC0PRUL",
	"#L0QVJG2",
	"#P2VYUGC",
	"#2JQR8VL",
	"#LCYUPRQ",
	"#Q8GLVJP",
	"#8R0YUVC",
	"#VJ2QPGL",
	"#2PVYUCR",
	"#9LURQG8",
	"#YGC0PVL",
	"#LQVJ8G2",
	"#PVCYUG2",
	"#2JRQV8L",
	"#LYUCPRQ",
	"#QGLVJ8P",
	"#8YUCV0R",
	"#VJQPG2L",
	"#2VUCPRQ",
	"#9URQGL8",
	"#YGCP0VL",
	"#LVJQG28",
}

// GetTopPlayerTags returns a copy of the top player tags list
// to prevent external modifications
func GetTopPlayerTags() []string {
	tags := make([]string, len(TopPlayerTags))
	copy(tags, TopPlayerTags)
	return tags
}

// GetTopPlayerCount returns the number of tracked players
func GetTopPlayerCount() int {
	return len(TopPlayerTags)
}
