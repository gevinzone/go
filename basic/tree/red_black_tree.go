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

// ReplaceBy 用节点r取代节点n
func (n *Node) ReplaceBy(r *Node) error {
	r.Parent = n.Parent
	r.Left = n.Left
	r.Right = n.Right
	if n.Parent == nil {
		return nil
	}
	if n.Parent.Left == n {
		n.Parent.Left = r
	} else {
		n.Parent.Right = r
	}
	return nil
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
		node.Left, node.Right = r.blackLeaf, r.blackLeaf
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

// insertFixCase3 关注节点是 n，它的叔叔节点 u 是黑色，关注节点 n 是其父节点 p 的左子节点
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

func (r *RedBlackTree) Delete(val int) error {
	var (
		node *Node
		err  error
	)
	if node, err = r.Search(val); err != nil {
		if err == ErrNodeNotFound {
			return nil
		}
		return err
	}
	return r.delete(node)
}

// delete 删除红黑树的节点，主要逻辑是，分2步完成删除：
// 1. 删除当前节点，并做初步调整，确保完成这一步后，树上的每个节点，依然满足该节点到其每个可达叶子节点的路径上，黑色节点数量相同
// 2. 在上一步基础上做进一步进行调整，保证每条路径上，不存在相邻红色节点
// 其中，第一步要分3个case进行处理，第二步分4个case进行处理
func (r *RedBlackTree) delete(node *Node) error {
	if r.deleteCase1(node) || r.deleteCase2(node) || r.deleteCase3(node) {
		return nil
	}
	return ErrNodeDeletingFailed
}

// deleteCase1 要删除的节点是 n，它只有一个子节点 c
// 算法逻辑是：
// 1. 删除节点 n，并且把节点 c 替换到节点 n 的位置，这一部分操作跟普通的二叉查找树的删除操作一样；
// 2. 节点 n 只能是黑色，节点 c 也只能是红色，其他情况均不符合红黑树的定义。这种情况下，我们把节点 c 改为黑色；
// 3. 调整结束，不需要进行二次调整。
func (r *RedBlackTree) deleteCase1(n *Node) bool {
	if n.Left != r.blackLeaf && n.Right != r.blackLeaf {
		return false
	}
	var c *Node
	if c = n.Left; c == r.blackLeaf {
		c = n.Right
	}

	if n.Parent.Left == n {
		n.Parent.Left = c
	} else {
		n.Parent.Right = c
	}
	c.SetBlack()
	//n = nil
	return true
}

// deleteCase2 要删除的节点 n 有两个非空子节点，并且它的后继节点就是节点 n 的右子节点 c
// 算法逻辑是：
// 1. 如果节点 n 的后继节点就是右子节点 c，那右子节点 c 肯定没有左子树。我们把节点 n 删除，并且将节点 c 替换到节点 n 的位置；
// 2. 然后把节点 c 的颜色设置为跟节点 n 相同的颜色；
// 3. 如果节点 c 是黑色，为了不违反红黑树的最后一条定义，我们给节点 c 的右子节点 d 多加一个黑色，这个时候节点 d 就成了“红 - 黑”或者“黑 - 黑”；
// 4. 这个时候，关注节点变成了节点 d，第二步的调整操作就会针对关注节点来做
func (r *RedBlackTree) deleteCase2(n *Node) bool {
	if n.Left == r.blackLeaf || n.Right == r.blackLeaf || n.Right != r.successor(n) {
		return false
	}
	c := n.Right
	c.Color = n.Color
	c.Left = n.Left
	c.Parent = n.Parent
	if n.Parent.Left == n {
		n.Parent.Left = c
	} else {
		n.Parent.Right = c
	}
	d := c.Right
	return r.deleteFix(d)
}

// deleteCase3 要删除的是节点 n，它有两个非空子节点，并且节点 n 的后继节点不是右子节点
// 算法逻辑是：
// 1. 找到后继节点 d，并将它删除，删除后继节点 d 的过程参照 CASE  1；
// 2. 将节点 n 替换成后继节点 d；
// 3. 把节点 d 的颜色设置为跟节点 n 相同的颜色；
// 4. 如果节点 d 是黑色，为了不违反红黑树的最后一条定义，我们给节点 d 的右子节点 c 多加一个黑色，这个时候节点 c 就成了“红 - 黑”或者“黑 - 黑”；
// 5. 这个时候，关注节点变成了节点 c，第二步的调整操作就会针对关注节点来做
func (r *RedBlackTree) deleteCase3(n *Node) bool {
	d := r.successor(n)
	if n.Left == r.blackLeaf || n.Right == r.blackLeaf || n.Right == d {
		return false
	}
	r.deleteCase1(d)
	if err := n.ReplaceBy(d); err != nil {
		return false
	}
	d.Color = n.Color
	c := d.Right
	return r.deleteFix(c)
}

// deleteFix 实现删除节点后的动态调整平衡
// 要分4种情况实现动态调整，目标是保证每条路径上，红色节点不相邻
func (r *RedBlackTree) deleteFix(n *Node) bool {
	if r.deleteFixCase1(n) || r.deleteFixCase2(n) || r.deleteFixCase3(n) || r.deleteFixCase4(n) {
		return true
	}
	return false
}

// deleteFixCase1 关注节点是 n，它的兄弟节点 b 是红色
// 调整算法为：
// 1. 围绕关注节点 n 的父节点 p 左旋；
// 2. 左旋后，b 成为 p 的父节点，将关注节点 n 的父节点 p 和祖父节点 b 交换颜色；
// 3. 关注节点不变，继续从四种情况中选择适合的规则来调整
func (r *RedBlackTree) deleteFixCase1(n *Node) bool {
	b := r.brother(n)
	if b == nil || b.IsBlack() {
		return false
	}
	// 既然有兄弟节点，p 一定不为空
	p := r.parent(n)
	r.leftRotate(p)
	// 左旋之后，p一定有父节点，即n一定有祖父节点
	p.Color, b.Color = b.Color, p.Color
	return r.deleteFix(n)
}

// deleteFixCase2 关注节点是 n，它的兄弟节点 b 是黑色的，并且节点 b 的左右子节点 d、e 都是黑色
// 调整算法为：
// 1. 将关注节点 n 的兄弟节点 b 的颜色变成红色；
// 2. 关注节点从 n 变成其父节点 p，继续从四种情况中选择符合的规则来调整。
func (r *RedBlackTree) deleteFixCase2(n *Node) bool {
	b := r.brother(n)
	if b == nil || b.IsRed() || b.Left.IsRed() || b.Right.IsRed() {
		return false
	}
	b.SetRed()
	// n有兄弟节点，则p一定不为空
	p := r.parent(n)
	return r.deleteFix(p)
}

// deleteFixCase3 关注节点是 n，它的兄弟节点 b 是黑色，b 的左子节点 d 是红色，b 的右子节点 e 是黑色
// 调整算法如下：
// 1. 围绕关注节点 n 的兄弟节点 b 右旋；
// 2. 节点 b 和节点 d 交换颜色；
// 3. 关注节点不变，跳转到 CASE 4，继续调整。
func (r *RedBlackTree) deleteFixCase3(n *Node) bool {
	b := r.brother(n)
	if b == nil || b.IsRed() || b.Left.IsBlack() || b.Right.IsRed() {
		return false
	}
	d := b.Left
	r.rightRotate(b)
	b.Color, d.Color = d.Color, b.Color
	return r.deleteFixCase4(n)
}

// deleteFixCase4 关注节点 n 的兄弟节点 b 是黑色的，并且 b 的右子节点是红色
// 调整算法如下：
// 1. 围绕关注节点 a 的父节点 p 左旋；
// 2. 将关注节点 a 的兄弟节点 b 的颜色，设置为关注节点 a 的父节点 p 的颜色；
// 3. 将关注节点 a 的父节点 p 的颜色设置为黑色；
// 4. 获取关注节点 a 的叔叔节点 e ，并将 e 设置为黑色；
// 5. 调整结束。
func (r *RedBlackTree) deleteFixCase4(n *Node) bool {
	b := r.brother(n)
	if b == nil || b.IsRed() || b.Right.IsBlack() {
		return false
	}
	p := r.parent(n)
	r.leftRotate(p)
	b.Color = p.Color
	p.SetBlack()
	e := r.brother(p)
	e.SetBlack()
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

// successor 找到删除n时的后继节点
// 由于n是红黑树的节点，符合二叉搜索树的要求，则查找后继节点的逻辑为：
// 1. 若n只有一个子节点，则该子节点即为后继节点
// 2. 若n有2个子节点，则后继节点为右子树的最小节点
// 		2.1 若n.Right.Left==nil，则Successor=n.Right
// 		2.2 若n.Right.Left!=nil，则Successor为n.Right.Left这个子树的最左节点
func (r *RedBlackTree) successor(n *Node) *Node {
	if n.Left == r.blackLeaf && n.Right == r.blackLeaf {
		return nil
	}
	if n.Left == r.blackLeaf || n.Right == r.blackLeaf {
		if n.Left == r.blackLeaf {
			return n.Right
		} else {
			return n.Left
		}
	}
	right := n.Right
	if right.Left == r.blackLeaf {
		return right
	}
	var p *Node
	for right != r.blackLeaf {
		p = right
		right = right.Left
	}
	return p
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
