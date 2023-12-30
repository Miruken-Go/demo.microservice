package credential

type Config struct {
	Clients []struct{
		Id     string
		Secret string
	}
}
