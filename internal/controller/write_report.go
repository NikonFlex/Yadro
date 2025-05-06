package controller

func (i implementation) WriteReport() error {
	i.logger.Debug("Writing report")
	err := i.competitorUseCase.WriteReport()
	if err != nil {
		return err
	}
	return nil
}
