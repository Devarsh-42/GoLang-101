package calculator

import "testing"

func TestAdd(t *testing.T) {
tests := []struct {
name string
a    int
b    int
want int
}{
{"positive numbers", 5, 3, 8},
{"negative numbers", -5, -3, -8},
{"mixed signs", 5, -3, 2},
{"with zero", 5, 0, 5},
{"both zero", 0, 0, 0},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := Add(tt.a, tt.b)
if got != tt.want {
t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
}
})
}
}

func TestSubtract(t *testing.T) {
tests := []struct {
name string
a    int
b    int
want int
}{
{"positive numbers", 10, 3, 7},
{"negative numbers", -10, -3, -7},
{"mixed signs", 10, -3, 13},
{"with zero", 10, 0, 10},
{"result negative", 3, 10, -7},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := Subtract(tt.a, tt.b)
if got != tt.want {
t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
}
})
}
}

func TestMultiply(t *testing.T) {
tests := []struct {
name string
a    int
b    int
want int
}{
{"positive numbers", 5, 3, 15},
{"negative numbers", -5, -3, 15},
{"mixed signs", 5, -3, -15},
{"with zero", 5, 0, 0},
{"with one", 5, 1, 5},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got := Multiply(tt.a, tt.b)
if got != tt.want {
t.Errorf("Multiply(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
}
})
}
}

func TestDivide(t *testing.T) {
tests := []struct {
name    string
a       int
b       int
want    int
wantErr bool
}{
{"positive numbers", 10, 2, 5, false},
{"negative dividend", -10, 2, -5, false},
{"negative divisor", 10, -2, -5, false},
{"both negative", -10, -2, 5, false},
{"division by zero", 10, 0, 0, true},
{"zero dividend", 0, 5, 0, false},
{"truncated result", 10, 3, 3, false},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
got, err := Divide(tt.a, tt.b)

if tt.wantErr {
if err == nil {
t.Errorf("Divide(%d, %d) expected error, got nil", tt.a, tt.b)
}
return
}

if err != nil {
t.Errorf("Divide(%d, %d) unexpected error: %v", tt.a, tt.b, err)
return
}

if got != tt.want {
t.Errorf("Divide(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
}
})
}
}
