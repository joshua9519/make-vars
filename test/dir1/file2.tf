resource "test" "name" {
  name = var.var1
  test = var.var2
  test2 = "${var.var3}-test"
}
