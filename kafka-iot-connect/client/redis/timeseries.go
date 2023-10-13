package redis

import (
	"context"
)

func (r *RedisTimeSeriesClient) XaddTrim(ctx context.Context, key, threshold, id, field, value string) {
	r.Client.Do(ctx, r.Client.B().Xadd().Key(key).Maxlen().Exact().Threshold(threshold).Id(id).FieldValue().FieldValue(field, value).Build())
}
