package clips

import (
	"testing"

	"gotest.tools/assert"
)

func TestFunctionsEnv(t *testing.T) {
	t.Run("List Functions", func(t *testing.T) {
		env := CreateEnvironment()
		defer env.Delete()

		err := env.Build(`(deffunction foo (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)
		err = env.Build(`(deffunction bar (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)

		funcs := env.Functions()
		assert.Equal(t, len(funcs), 2)
		assert.Equal(t, funcs[0].Name(), "foo")
		assert.Equal(t, funcs[1].Name(), "bar")
	})

	t.Run("Find Function", func(t *testing.T) {
		env := CreateEnvironment()
		defer env.Delete()

		err := env.Build(`(deffunction foo (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)
		err = env.Build(`(deffunction bar (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)

		ftion, err := env.FindFunction("foo")
		assert.NilError(t, err)
		assert.Equal(t, ftion.Name(), "foo")

		_, err = env.FindFunction("baz")
		assert.ErrorContains(t, err, "not found")
	})
}

func TestFunctions(t *testing.T) {
	t.Run("Function basics", func(t *testing.T) {
		env := CreateEnvironment()
		defer env.Delete()

		err := env.Build(`(deffunction foo (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)

		ftion, err := env.FindFunction("foo")
		assert.NilError(t, err)
		assert.Equal(t, ftion.Name(), "foo")
		assert.Equal(t, ftion.String(), `(deffunction MAIN::foo
   (?a ?b)
   (+ ?a ?b))`)
	})

	t.Run("Function equals", func(t *testing.T) {
		env := CreateEnvironment()
		defer env.Delete()

		err := env.Build(`(deffunction foo (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)
		err = env.Build(`(deffunction bar (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)

		ftion, err := env.FindFunction("foo")
		assert.NilError(t, err)
		ftion2, err := env.FindFunction("foo")
		assert.NilError(t, err)
		assert.Assert(t, ftion.Equals(ftion2))
		ftion2, err = env.FindFunction("bar")
		assert.NilError(t, err)
		assert.Assert(t, !ftion.Equals(ftion2))
	})

	t.Run("Function call", func(t *testing.T) {
		env := CreateEnvironment()
		defer env.Delete()

		err := env.Build(`(deffunction foo (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)

		ftion, err := env.FindFunction("foo")
		assert.NilError(t, err)

		ret, err := ftion.Call("1 2")
		assert.NilError(t, err)
		assert.Equal(t, ret, int64(3))
	})

	t.Run("Module", func(t *testing.T) {
		env := CreateEnvironment()
		defer env.Delete()

		err := env.Build(`(deffunction foo (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)

		ftion, err := env.FindFunction("foo")
		assert.NilError(t, err)

		mod := ftion.Module()
		assert.Equal(t, mod.Name(), "MAIN")
	})

	t.Run("Deletable", func(t *testing.T) {
		env := CreateEnvironment()
		defer env.Delete()

		err := env.Build(`(deffunction foo (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)

		ftion, err := env.FindFunction("foo")
		assert.NilError(t, err)

		assert.Assert(t, ftion.Deletable())

		err = env.Build(`(defrule fooref => (foo 1 2))`)
		assert.NilError(t, err)

		assert.Assert(t, !ftion.Deletable())
	})

	t.Run("Undefine", func(t *testing.T) {
		env := CreateEnvironment()
		defer env.Delete()

		err := env.Build(`(deffunction foo (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)

		ftion, err := env.FindFunction("foo")
		assert.NilError(t, err)
		err = env.Build(`(defrule fooref => (foo 1 2))`)
		assert.NilError(t, err)

		err = ftion.Undefine()
		assert.ErrorContains(t, err, "Unable")

		_, err = env.Eval(`(undefrule fooref)`)
		assert.NilError(t, err)
		err = ftion.Undefine()
		assert.NilError(t, err)
	})

	t.Run("Watch", func(t *testing.T) {
		env := CreateEnvironment()
		defer env.Delete()

		err := env.Build(`(deffunction foo (?a ?b) (+ ?a ?b))`)
		assert.NilError(t, err)

		ftion, err := env.FindFunction("foo")
		assert.NilError(t, err)

		assert.Assert(t, !ftion.Watched())
		ftion.Watch(true)
		assert.Assert(t, ftion.Watched())
	})
}
