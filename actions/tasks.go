package actions

import (
	"net/http"
	"strings"

	"todo/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// TasksIndex muestra la lista de tareas y el formulario de creación.
func TasksIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	return renderIndex(c, tx, http.StatusOK, &models.Task{}, validate.NewErrors(), nil, validate.NewErrors())
}

// TasksCreate crea una tarea nueva a partir del formulario.
func TasksCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	task := &models.Task{}
	if err := c.Bind(task); err != nil {
		return err
	}
	task.Title = strings.TrimSpace(task.Title)

	verrs, err := tx.ValidateAndCreate(task)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return renderIndex(c, tx, http.StatusUnprocessableEntity, task, verrs, nil, validate.NewErrors())
	}

	c.Flash().Add("success", "Tarea creada.")
	return c.Redirect(http.StatusSeeOther, "/")
}

// TasksEdit muestra la lista con el formulario de edición de una tarea.
func TasksEdit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	task := &models.Task{}
	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return renderIndex(c, tx, http.StatusOK, &models.Task{}, validate.NewErrors(), task, validate.NewErrors())
}

// TasksUpdate cambia el título de una tarea.
func TasksUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	task := &models.Task{}
	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	form := &models.Task{}
	if err := c.Bind(form); err != nil {
		return err
	}
	task.Title = strings.TrimSpace(form.Title)

	verrs, err := tx.ValidateAndUpdate(task)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return renderIndex(c, tx, http.StatusUnprocessableEntity, &models.Task{}, validate.NewErrors(), task, verrs)
	}

	c.Flash().Add("success", "Tarea actualizada.")
	return c.Redirect(http.StatusSeeOther, "/")
}

// renderIndex pinta la lista de tareas. editing es la tarea en edición, o nil.
func renderIndex(c buffalo.Context, tx *pop.Connection, status int, task *models.Task, verrs *validate.Errors, editing *models.Task, editErrs *validate.Errors) error {
	tasks := models.Tasks{}
	if err := tx.Order("created_at asc, id asc").All(&tasks); err != nil {
		return err
	}

	editingID := ""
	editingTitle := ""
	if editing != nil {
		editingID = editing.ID.String()
		editingTitle = editing.Title
	}

	c.Set("tasks", tasks)
	c.Set("task", task)
	c.Set("errors", verrs)
	c.Set("editingID", editingID)
	c.Set("editingTitle", editingTitle)
	c.Set("editErrors", editErrs)
	return c.Render(status, r.HTML("tasks/index.plush.html"))
}

// TasksToggle alterna una tarea entre completada y pendiente.
func TasksToggle(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	task := &models.Task{}
	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	task.Completed = !task.Completed
	if err := tx.Update(task); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

// TasksDestroy elimina una tarea.
func TasksDestroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)

	task := &models.Task{}
	if err := tx.Find(task, c.Param("task_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(task); err != nil {
		return err
	}

	c.Flash().Add("success", "Tarea eliminada.")
	return c.Redirect(http.StatusSeeOther, "/")
}
