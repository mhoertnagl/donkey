package aegis

type color struct {
  r uint8
  g uint8
  b uint8
}

func Color(r, g, b uint8) *color {
  return &color{r: r, g: g, b: b};
}

func Grey(v uint8) *color {
  return Color(v, v, v);
}

func White() *color {
  return Grey(255);
}
