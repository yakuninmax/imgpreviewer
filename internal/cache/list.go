package cache

type list struct {
	len   int
	front *listItem
	back  *listItem
}
type listItem struct {
	value image
	next  *listItem
	prev  *listItem
}
