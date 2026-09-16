package httpapi

import (
	"github.com/eduardoaugustolb/versum/api/internal/catalog/application/queries"
	"github.com/eduardoaugustolb/versum/api/internal/health"
	identityAccessPolicy "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/policy"
	identityAccessPorts "github.com/eduardoaugustolb/versum/api/internal/identityaccess/application/ports"
	outboxPorts "github.com/eduardoaugustolb/versum/api/internal/outboxevent/application/ports"
	clockPorts "github.com/eduardoaugustolb/versum/api/internal/ports/clock"
	"github.com/eduardoaugustolb/versum/api/internal/ports/dbexec"
	idGeneratorPorts "github.com/eduardoaugustolb/versum/api/internal/ports/id"
)

type Dependencies struct {
	Health         health.CheckHealth
	Catalog        CatalogDependencies
	IdentityAccess IdentityAccessDependencies
}

type CatalogDependencies struct {
	ListBooks  *queries.ListBooks
	GetChapter *queries.GetChapter
}

type IdentityAccessDependencies struct {
	DbExecutor             dbexec.Executor
	EmailProtector         identityAccessPorts.EmailProtector
	LoginTokenGenerator    identityAccessPorts.LoginTokenGenerator
	LoginTokenHasher       identityAccessPorts.LoginTokenHasher
	IdGenerator            idGeneratorPorts.IDGenerator
	Clock                  clockPorts.Clock
	MagicLinkPolicy        identityAccessPolicy.MagicLinkPolicy
	OutboxPayloadProtector outboxPorts.PayloadProtector
}
