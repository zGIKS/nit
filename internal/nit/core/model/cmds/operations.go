package cmds

import (
	"github.com/zGIKS/nit/internal/nit/app"
	g "github.com/zGIKS/nit/internal/nit/git"
)

func ExecOperation(svc g.Service, op app.Operation) (string, error) {
	switch op.Kind {
	case app.OpStagePath:
		return svc.StagePath(op.Path)
	case app.OpUnstagePath:
		return svc.UnstagePath(op.Path)
	case app.OpStageAll:
		return svc.StageAll()
	case app.OpUnstageAll:
		return svc.UnstageAll()
	case app.OpDiscardAll:
		return svc.DiscardAll()
	case app.OpCommit:
		return svc.CommitWithOptions(op.Message, g.CommitOptions{
			All:     op.CommitAll,
			Amend:   op.CommitAmend,
			Signoff: op.CommitSignoff,
		})
	case app.OpPull:
		return svc.Pull()
	case app.OpFetch:
		return svc.Fetch()
	case app.OpPush:
		return svc.Push()
	case app.OpUndoLastCommit:
		return svc.UndoLastCommit()
	case app.OpAbortRebase:
		return svc.AbortRebase()
	default:
		return "", nil
	}
}
