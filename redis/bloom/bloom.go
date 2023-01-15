package bloom

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
)

func bloomExample(ctx context.Context, client redis.UniversalClient) error {
	if err := initData(ctx, client); err != nil {
		return err
	}
	if err := dataExistsExample(ctx, client); err != nil {
		return err
	}
	if err := multiInsertExample(ctx, client); err != nil {
		return err
	}

	return nil
}

func initData(ctx context.Context, client redis.UniversalClient) error {
	inserted, err := client.Do(ctx, "BF.ADD", "bf_key", "item0").Bool()
	if err != nil {
		return err
	}
	if inserted {
		fmt.Println("item0 is inserted")
		return nil
	}
	fmt.Println("item0 is already existed")
	return nil
}

func dataExistsExample(ctx context.Context, client redis.UniversalClient) error {
	data := []string{"item0", "item1"}
	for _, d := range data {
		exist, err := client.Do(ctx, "BF.EXISTS", "bf_key", d).Bool()
		if err != nil {
			return err
		}
		if exist {
			fmt.Println(d, " exists")
		} else {
			fmt.Println(d, " does not exist")
		}
	}
	return nil
}

func multiInsertExample(ctx context.Context, client redis.UniversalClient) error {
	values := []any{"BF.MADD", "bf_key", "item1", "item2", "item3"}
	booleans, err := client.Do(ctx, values...).BoolSlice()
	if err != nil {
		return err
	}
	fmt.Println("adding multiple items:", booleans)
	return nil
}
