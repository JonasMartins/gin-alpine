// Package utils
package utils

import "time"

type ContextString string

var AuthUserIDField = ContextString("middleware.auth.user_id")
var AuthUser = ContextString("middleware.auth.user")

const (
	CookieAccessToken              = "_cookie_at_h7tKlm6DDz"
	CookieRefreshToken             = "_cookie_rf_vTlFTGnBZ2"
	BearerTokenRegex               = `(?i)bearer\s+([a-zA-Z0-9\-_]+)`
	Minute                         = 60
	Hour                           = Minute * 60
	TaskSendingEmailType           = "email:deliver"
	TaskDCategorias                = "TaskDCategorias"
	TaskFiliais                    = "TASK_FILIAIS"
	TaskFamilias                   = "TASK_FAMILIAS"
	TaskDClientes                  = "TaskDClientes"
	TaskFContasPagar               = "TASK_FCONTAS_PAGAR"
	TaskFContasReceber             = "TASK_FCONTAS_RECEBER"
	TaskDProdutoFornecedor         = "TASK_DPRODUTO_FORNECEDOR"
	TaskDProduto                   = "TASK_DPRODUTO"
	TaskFPedido                    = "TASK_FPEDIDO"
	TaskDPedidoAutorizado          = "TASK_DPEDIDO_AUTORIZADO"
	TaskUpdateCMCDPedidoAutorizado = "TASK_UPDATE_CMC_DPEDIDO_AUTORIZADO"
	TaskFPedidoCancelado           = "TASK_FPEDIDO_CANCELADO"
	TaskFPedidosCompras            = "TASK_FPEDIDOS_COMPRAS"
	TaskDUsuario                   = "TASK_DUSUARIO"
	TaskDDepartamentos             = "TASK_DDEPARTAMENTOS"
	TaskDMovEstoque                = "TASK_DMOV_ESTOQUE"
	TaskFMovimentacaoReceber       = "TASK_FMOVIMENTACAO_RECEBER"
	TaskDPosicaoEstoque            = "TASK_DPOSICAO_ESTOQUE"
	TaskFFluxoCaixa                = "TASK_FFLUXO_CAIXA"
	TaskDContasCorrentes           = "TASK_DCONTAS_CORRENTES"
	TaskAddGolampExecucao          = "TASK_ADD_GOLAMP_EXECUCAO"
	TaskUpdateGolampExecucao       = "TASK_UPDATE_GOLAMP_EXECUCAO"
	JobRotateLogs                  = "ROTATE_LOGS"
	JobRunCarga                    = "RUN_CARGA"
	JobSendingClientsEmails        = "SENDING_CLIENTS_EMAILS"
	SendResetPasswordTask          = "SEND_RESET_PASSWORD_TASK"
	SendLoadFinishedTask           = "SEND_LOAD_FINISHED_TASK"
	EnqueueTask                    = "ENQUEUE_TASK"
	JobStarted                     = "JOB_STARTED"
	JobFinished                    = "JobFinished"
	TaskSendEmail                  = "TASK_SENDING_EMAIL"
	LinkResetPasswordExpiration    = 3 * time.Hour
	LongCachePastData              = 720 * time.Hour
	ShortCacheActualData           = 24 * time.Hour
	CronEveryMinute                = "* * * * *"
	DBInsert                       = "INSERT_INTO_DATABASE"

	// Logs
	RouteRefreshToken                 = "ROUTE_REFRESH_TOKEN"
	RouteLogout                       = "ROUTE_LOGOUT"
	RouteLogin                        = "ROUTE_LOGIN"
	RouteGetLink                      = "ROUTE_GET_LINK"
	RouteCreateResetPass              = "ROUTE_CREATE_RESET_PASS"
	RouteCreateUser                   = "ROUTE_CREATE_USER"
	RouteUpdateUser                   = "ROUTE_UPDATE_USER"
	RouteUpdateUserAdmin              = "ROUTE_UPDATE_USER_ADMIN"
	RouteUpdateUserManager            = "ROUTE_UPDATE_USER_MANAGER"
	RouteGetUsers                     = "ROUTE_GET_USERS"
	RouteGetUser                      = "ROUTE_GET_USER"
	RouteFindUser                     = "ROUTE_FIND_USER"
	RouteResetPasswordFromLink        = "ROUTE_RESET_PASSWORD_FROM_LINK"
	RouteGetSellers                   = "ROUTE_GET_SELLERS"
	RouteGetPedidosByDates            = "ROUTE_GET_PEDIDOS_BY_DATES"
	RouteGetTotalOrdersByDates        = "ROUTE_GET_TOTAL_ORDERS_BY_DATES"
	RouteGetTotalOrdersByDatesSeller  = "ROUTE_GET_TOTAL_ORDERS_BY_DATES_AND_SELLER"
	RouteGetAllOrdersByDate           = "ROUTE_GET_ALL_ORDERS_BY_DATE"
	RouteGetSumOrdersFromPrevMonths   = "ROUTE_GET_SUM_ORDERS_FROM_PREVIOUS_MONTHS"
	RouteGetBranches                  = "ROUTE_GET_BRANCHES"
	RouteGetSuppliersDetailsFromMonth = "ROUTE_GET_SUPPLIERS_DETAILS_FROM_MONTH"
	// Cache
	UsersPrefixCache    = "users"
	FPedidosPrefixCache = "fpedidos"

	// Errors
	EmailAlreadyExists       = "EMAIL_ALREADY_EXISTS"
	WrongPassword            = "WRONG PASSWORD"
	NotFound                 = "NOT FOUND"
	ServerError              = "SERVER ERROR"
	TokenGenerateError       = "FAIL TO GENERATE TOKEN"
	MigrationNoChange        = "NO CHANGE"
	TokenExpired             = "TOKEN EXPIRED"
	FailToParseCustomError   = "FAIL TO PARSE CUSTOM ERROR"
	InvalidToken             = "INVALID TOKEN"
	FailToReadBodyResponse   = "FAIL TO READ BODY RESPONSE"
	FailToReadBodyRequest    = "FAIL TO READ BODY REQUEST"
	FailToDecodeBodyResponse = "FAIL TO DECODE BODY RESPONSE"
	APIKeyRequired           = "VALID API KEY REQUIRED"
)
