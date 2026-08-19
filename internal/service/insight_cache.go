package service

type InsightCache struct {
	items    map[string]ExpeditionInsight
	order    []string
	capacity int
}

func NewInsightCache() *InsightCache {
	return &InsightCache{
		items:    make(map[string]ExpeditionInsight),
		order:    make([]string, 0, 8),
		capacity: 16,
	}
}

func (c *InsightCache) Get(key string) (ExpeditionInsight, bool) {
	value, ok := c.items[key]
	if !ok {
		return ExpeditionInsight{}, false
	}
	c.touch(key)
	return value, true
}

func (c *InsightCache) Put(key string, value ExpeditionInsight) {
	if _, exists := c.items[key]; exists {
		c.items[key] = value
		c.touch(key)
		return
	}
	if len(c.order) >= c.capacity {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.items, oldest)
	}
	c.items[key] = value
	c.order = append(c.order, key)
}

func (c *InsightCache) Invalidate(key string) {
	if _, exists := c.items[key]; !exists {
		return
	}
	c.remove(key)
}

func (c *InsightCache) remove(key string) {
	delete(c.items, key)
	for index, item := range c.order {
		if item != key {
			continue
		}
		c.order = append(c.order[:index], c.order[index+1:]...)
		return
	}
}

func (c *InsightCache) touch(key string) {
	for index, item := range c.order {
		if item == key {
			c.order = append(c.order[:index], c.order[index+1:]...)
			break
		}
	}
	c.order = append(c.order, key)
}
