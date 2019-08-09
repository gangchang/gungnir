package gungnir

type beforeExecFn func(Ctx) bool

type afterExecFn func(Ctx) bool
