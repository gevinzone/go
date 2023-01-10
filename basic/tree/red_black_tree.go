package tree

type Colors int

const (
	BLACK Colors = iota
	RED
)

type Node struct {
	Color  Colors
	Parent *Node
	Left   *Node
	Right  *Node
	val    int
}

func (n *Node) IsBlack() bool {
	return n.Color == BLACK
}

func (n *Node) IsRed() bool {
	return n.Color == RED
}

func (n *Node) SetBlack() {
	n.Color = BLACK
}

func (n *Node) SetRed() {
	n.Color = RED
}

type RedBlackTree struct {
	Root      *Node
	blackLeaf *Node
}

type RBTOption func(tree *RedBlackTree)

func NewRedBlackTree(opts ...RBTOption) *RedBlackTree {
	res := &RedBlackTree{
		blackLeaf: &Node{Color: BLACK},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (r *RedBlackTree) Find(val int) (*Node, error) {
	return r.find(r.Root, val)
}

func (r *RedBlackTree) find(node *Node, val int) (*Node, error) {
	if node == nil {
		return nil, ErrNodeNotFound
	}
	if node.val == val {
		return r.Root, nil
	}
	if node.val > val {
		return r.find(node.Left, val)
	}
	return r.find(node.Right, val)
}

func (r *RedBlackTree) Search(val int) (*Node, error) {
	if r.Root == nil {
		return nil, ErrNodeNotFound
	}
	n := r.Root
	for n != r.blackLeaf {
		if n.val == val {
			return n, nil
		}
		if n.val > val {
			n = n.Left
			continue
		}
		if n.val < val {
			n = n.Right
		}
	}
	return nil, ErrNodeNotFound
}

func (r *RedBlackTree) Insert(val int) error {
	node := &Node{val: val, Color: RED}
	err := r.insert(node)
	if err != nil {
		return err
	}
	return r.insertFix(node)
}

func (r *RedBlackTree) insert(node *Node) error {
	if r.Root == nil {
		r.Root = node
		return nil
	}
	n := r.Root
	var p *Node
	for n != r.blackLeaf {
		p = n
		if n.val == node.val {
			return ErrNodeExisting
		}
		if n.val > node.val {
			n = n.Left
			continue
		}
		n = n.Right
	}
	if p.val > node.val {
		p.Left = node
	} else {
		p.Right = node
	}
	node.Parent = p
	node.Left, node.Right = r.blackLeaf, r.blackLeaf
	return nil
}

// insertFix 实现插入新节点后的动态调整平衡
// 红黑树插入新节点的动态平衡算法如下：
// 1. 插入node为root时，直接变为黑色即可；
// 2. 插入节点的父节点为黑时，已经符合要求，不用调整
// 3. 插入节点的父节点为红时，需要动态调整，分以下3个case分别处理：

func (r *RedBlackTree) insertFix(node *Node) error {
	if node == r.Root {
		node.SetBlack()
		return nil
	}
	if node.Parent.IsBlack() {
		return nil
	}
	return r.insertFixByCase(node)
}

func (r *RedBlackTree) insertFixByCase(n *Node) error {
	if r.insertFixCase1(n) || r.insertFixCase2(n) || r.insertFixCase3(n) {
		return nil
	}

	return ErrTreeBalancingFailed
}

func (r *RedBlackTree) getRelativeNodes(node *Node) (parent *Node, uncle *Node, grandParent *Node) {
	parent = r.parent(node)
	uncle = r.brother(parent)
	grandParent = r.parent(parent)
	return
}

// insertFixCase1 关注节点是 n，它的叔叔节点 u 是红色
// 调整算法为：
// 1. 将关注节点 n 的父节点 p、叔叔节点 u 的颜色都设置成黑色；
// 2. 将关注节点 n 的祖父节点 g 的颜色设置成红色；
// 3. 关注节点变成 n 的祖父节点 g；
// 4. 跳到 CASE 2 或者 CASE 3。
func (r *RedBlackTree) insertFixCase1(n *Node) bool {
	p, u, g := r.getRelativeNodes(n)
	if u.IsBlack() {
		return false
	}
	p.SetBlack()
	u.SetBlack()
	g.SetRed()
	if r.insertFixCase2(g) || r.insertFixCase3(g) {
		return true
	}
	return false
}

// insertFixCase2 关注节点是 n，它的叔叔节点 u 是黑色，关注节点 n 是其父节点 p 的右子节点
// 调整算法为：
// 1. 关注节点变成节点 n 的父节点 p；
// 2. 围绕新的关注节点 p 左旋；
// 3. 跳到 CASE 3。
func (r *RedBlackTree) insertFixCase2(n *Node) bool {
	p, u, _ := r.getRelativeNodes(n)
	if u.IsRed() || n != p.Right {
		return false
	}
	r.leftRotate(p)
	return r.insertFixCase3(p)
}

// insertFixCase3 关注节点是 a，它的叔叔节点 u 是黑色，关注节点 n 是其父节点 p 的左子节点
// 调整算法为：
// 1. 围绕关注节点 n 的祖父节点 g 右旋；
// 2. 将关注节点 n 的父节点 p、兄弟节点 g 的颜色互换(围绕g右旋后，n的原祖父节点g变为n的兄弟节点)。
// 3. 调整结束。
func (r *RedBlackTree) insertFixCase3(n *Node) bool {
	p, u, g := r.getRelativeNodes(n)
	if u.IsRed() || n != p.Left {
		return false
	}
	r.rightRotate(g)
	p.Color, g.Color = g.Color, p.Color
	return true
}

func (r *RedBlackTree) parent(node *Node) *Node {
	if node == nil {
		return nil
	}
	return node.Parent
}

func (r *RedBlackTree) brother(node *Node) *Node {
	if node == nil || node.Parent == nil {
		return nil
	}
	p := node.Parent
	if node == p.Left {
		return p.Right
	}
	return p.Left
}

func (r *RedBlackTree) leftRotate(n *Node) {
	if n == nil || n.Right == r.blackLeaf {
		return
	}
	p, x, y := r.parent(n), n, n.Right
	if p == nil {
		r.Root = y
	} else {
		switch x {
		case p.Left:
			p.Left = y
		case p.Right:
			p.Right = y
		}
	}
	x.Parent, y.Parent = y, p
	if y.Left != r.blackLeaf {
		y.Left.Parent = x
	}
	x.Right, y.Left = y.Left, x
}

func (r *RedBlackTree) rightRotate(n *Node) {
	if n == nil || n.Left == r.blackLeaf {
		return
	}
	p, x, y := r.parent(n), n, n.Left
	if p == nil {
		r.Root = y
	} else {
		switch x {
		case p.Left:
			p.Left = y
		case p.Right:
			p.Right = y
		}
	}
	x.Parent, y.Parent = y, p
	if y.Right != nil {
		y.Right.Parent = x
	}
	x.Left, y.Right = y.Right, x
}
