package controller

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type OpenFileDialogArgs struct {
	Title string
}

func (c *Controller) OpenFileDialog(args OpenFileDialogArgs) ([]string, error) {
	var filter []runtime.FileFilter
	for i := range c.appConfig.UI.FileDialog.FilterDisplay {
		filter = append(filter, runtime.FileFilter{
			DisplayName: c.appConfig.UI.FileDialog.FilterDisplay[i],
			Pattern:     c.appConfig.UI.FileDialog.FilterPattern[i],
		})
	}

	return runtime.OpenMultipleFilesDialog(c.ctx, runtime.OpenDialogOptions{
		Title:                      args.Title,
		Filters:                    filter,
		DefaultDirectory:           c.appConfig.UI.FileDialog.DefaultDirectory,
		ShowHiddenFiles:            c.appConfig.UI.FileDialog.ShowHiddenFiles,
		CanCreateDirectories:       c.appConfig.UI.FileDialog.CanCreateDirectories,
		ResolvesAliases:            c.appConfig.UI.FileDialog.ResolveAliases,
		TreatPackagesAsDirectories: c.appConfig.UI.FileDialog.TreatPackagesAsDirectories,
	})
}
