package processor

import (
	"context"
	"fmt"
)

type updateOperationsImpl struct {
	processor *processor
}

func newUpdateOperations(p *processor) updateOperations {
	return &updateOperationsImpl{processor: p}
}

func (u *updateOperationsImpl) downloadInstaller(_ context.Context, state *operationState) error {
	u.processor.appendLog(MsgDownloadingInstaller, ProductStackInstaller, state)
	return fmt.Errorf("not implemented")
}

func (u *updateOperationsImpl) updateInstaller(_ context.Context, state *operationState) error {
	u.processor.appendLog(MsgUpdatingInstaller, ProductStackInstaller, state)
	return fmt.Errorf("not implemented")
}

func (u *updateOperationsImpl) removeInstaller(_ context.Context, state *operationState) error {
	u.processor.appendLog(MsgRemovingInstaller, ProductStackInstaller, state)
	return fmt.Errorf("not implemented")
}

