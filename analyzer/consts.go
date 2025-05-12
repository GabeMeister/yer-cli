package analyzer

const PREV_YEAR = 2023
const CURR_YEAR = 2024

var MONTHS = [...]string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

var WEEK_DAYS = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

var HOURS = [...]string{
	"12am",
	"1am",
	"2am",
	"3am",
	"4am",
	"5am",
	"6am",
	"7am",
	"8am",
	"9am",
	"10am",
	"11am",
	"12pm",
	"1pm",
	"2pm",
	"3pm",
	"4pm",
	"5pm",
	"6pm",
	"7pm",
	"8pm",
	"9pm",
	"10pm",
	"11pm",
}

var SUPPORTED_FILE_EXTENSIONS = []string{
	"js", "mjs", "cjs", "jsx", "javascript", "py", "pyw", "pyc", "pyd", "pyo", "pyi", "pyz",
	"java", "class", "jar", "cs", "csx", "cpp", "cc", "cxx", "c++", "hpp", "hh", "hxx", "h++",
	"ts", "tsx", "mts", "cts", "php", "phtml", "php3", "php4", "php5", "php7", "phps", "php-s",
	"phar", "c", "h", "go", "mod", "rs", "rlib", "swift", "swiftdoc", "swiftmodule", "kt", "kts",
	"ktm", "rb", "rbw", "rake", "gemspec", "ru", "r", "rdata", "rds", "rmd", "dart", "scala",
	"sc", "sql", "ddl", "dml", "html", "htm", "xhtml", "css", "scss", "sass", "less", "m", "mm",
	"sh", "zsh", "fish", "pl", "pm", "t", "pod", "lua", "luac", "mat", "fig", "mlx", "mlapp",
	"groovy", "gvy", "gy", "gsh", "hs", "lhs", "asm", "s", "S", "clj", "cljs", "cljc", "edn",
	"ex", "exs", "heex", "jl", "fs", "fsi", "fsx", "fsscript", "vba", "bas", "cls", "frm",
	"ps1", "psm1", "psd1", "ps1xml", "pssc", "cdxml", "vb", "vbs", "cob", "cbl", "cpy", "f",
	"for", "f90", "f95", "f03", "f08", "cr", "d", "di", "erl", "hrl", "beam", "app", "escript",
	"coffee", "litcoffee", "pas", "pp", "dpk", "dpr", "dproj", "apex", "apxc", "apxt", "bash",
	"abap", "ada", "adb", "ads", "pro", "p", "lisp", "lsp", "l", "cl", "fasl", "sas", "sas7bdat",
	"ml", "mli", "cmx", "cmo", "scm", "ss", "sls", "sps", "hack", "hh", "sol", "rkt", "rktl",
	"rktd", "scrbl", "vhd", "vhdl", "v", "vh", "sv", "svh", "elm", "pls", "plb", "pck", "pks",
	"pkb", "inc", "wasm", "wat", "zig", "nim", "nims", "st", "as", "mxml", "logo", "lg", "tcl",
	"tk", "purs", "hx", "hxml", "raku", "rakumod", "rakudoc", "rakutest", "p6", "pl6", "pm6",
	"re", "rei", "vi", "lvproj", "lvlib", "fth", "4th", "forth", "j", "ijs", "ijt", "apl",
	"aplf", "aplo", "apln", "aplc", "xslt", "xsl", "ksh", "qs", "vala", "vapi", "sml", "sig",
	"fun", "bal", "res", "resi", "io", "red", "reds", "factor", "icl", "dcl", "dylan", "dyl",
	"intr", "als", "agda", "lagda", "idr", "lidr", "coq", "wl", "nb", "wls", "odin", "ls",
	"livescript", "moo", "e", "ace", "pike", "pmod", "jse", "yml", "yaml", "md",
}
