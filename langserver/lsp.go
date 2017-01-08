package langserver

type InitializeParams struct {
	ProcessID             int                `json:"processId,omitempty"`
	RootPath              string             `json:"rootPath,omitempty"`
	InitializationOptions InitializeOptions  `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities",omitempty`
	Trace                 string             `json:"trace,omitempty"`
}

type InitializeOptions struct {
}

type ClientCapabilities struct {
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities,omitempty"`
}

type TextDocumentSyncKind int

const (
	TDSKNone        TextDocumentSyncKind = 0
	TDSKFull                             = 1
	TDSKIncremental                      = 2
)

type ServerCapabilities struct {
	TextDocumentSync       TextDocumentSyncKind `json:"textDocumentSync,omitempty"`
	DocumentSymbolProvider bool                 `json:"documentSymbolProvider,omitempty"`
}

type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageId string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier
}

type SymbolKind int

const (
	SKFile        SymbolKind = 1
	SKModule      SymbolKind = 2
	SKNamespace   SymbolKind = 3
	SKPackage     SymbolKind = 4
	SKClass       SymbolKind = 5
	SKMethod      SymbolKind = 6
	SKProperty    SymbolKind = 7
	SKField       SymbolKind = 8
	SKConstructor SymbolKind = 9
	SKEnum        SymbolKind = 10
	SKInterface   SymbolKind = 11
	SKFunction    SymbolKind = 12
	SKVariable    SymbolKind = 13
	SKConstant    SymbolKind = 14
	SKString      SymbolKind = 15
	SKNumber      SymbolKind = 16
	SKBoolean     SymbolKind = 17
	SKArray       SymbolKind = 18
)

type SymbolInformation struct {
	Name     string     `json:"name"`
	Kind     SymbolKind `json:"kind"`
	Location Location   `json:"location"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}
