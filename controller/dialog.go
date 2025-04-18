package controller

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type OpenFileDialogArgs struct {
	Title string
}

func (c *Controller) OpenFileDialog(args OpenFileDialogArgs) ([]string, error) {
	var filter []runtime.FileFilter
	for i := range c.getConfig().UI.FileDialog.FilterDisplay {
		filter = append(filter, runtime.FileFilter{
			DisplayName: c.getConfig().UI.FileDialog.FilterDisplay[i],
			Pattern:     c.getConfig().UI.FileDialog.FilterPattern[i],
		})
	}

	return runtime.OpenMultipleFilesDialog(c.ctx, runtime.OpenDialogOptions{
		Title:                      args.Title,
		Filters:                    filter,
		DefaultDirectory:           c.getConfig().UI.FileDialog.DefaultDirectory,
		ShowHiddenFiles:            c.getConfig().UI.FileDialog.ShowHiddenFiles,
		CanCreateDirectories:       c.getConfig().UI.FileDialog.CanCreateDirectories,
		ResolvesAliases:            c.getConfig().UI.FileDialog.ResolveAliases,
		TreatPackagesAsDirectories: c.getConfig().UI.FileDialog.TreatPackagesAsDirectories,
	})
}
