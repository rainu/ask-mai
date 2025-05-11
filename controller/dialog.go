package controller

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type OpenFileDialogArgs struct {
	Title string
}

func (c *Controller) OpenFileDialog(args OpenFileDialogArgs) ([]string, error) {
	var filter []runtime.FileFilter
	for i := range c.getProfile().UI.FileDialog.FilterDisplay {
		filter = append(filter, runtime.FileFilter{
			DisplayName: c.getProfile().UI.FileDialog.FilterDisplay[i],
			Pattern:     c.getProfile().UI.FileDialog.FilterPattern[i],
		})
	}

	return runtime.OpenMultipleFilesDialog(c.ctx, runtime.OpenDialogOptions{
		Title:                      args.Title,
		Filters:                    filter,
		DefaultDirectory:           c.getProfile().UI.FileDialog.DefaultDirectory,
		ShowHiddenFiles:            *c.getProfile().UI.FileDialog.ShowHiddenFiles,
		CanCreateDirectories:       *c.getProfile().UI.FileDialog.CanCreateDirectories,
		ResolvesAliases:            *c.getProfile().UI.FileDialog.ResolveAliases,
		TreatPackagesAsDirectories: *c.getProfile().UI.FileDialog.TreatPackagesAsDirectories,
	})
}
