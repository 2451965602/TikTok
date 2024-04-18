namespace go video

include "model.thrift"


//视频流
struct FeedRequest{
    1:optional string latest_time(api.query="latest_time"),
    2:required i64 page_num (api.query="page_num"),
    3:required i64 page_size (api.query="page_size")
}

struct FeedResponse{
    1:model.BaseResp base,
    2:model.VideoList data,
}

//投稿
struct UploadRequest{
    1:required string title (api.form="title"),
//    2:required binary coverdata (api.form="coverdata"),
//    3:required binary videodata (api.form="videodata"),
    2:required string description (api.form="description"),
}

struct UploadResponse{
    1:model.BaseResp base,
}
//发布列表
struct UploadListRequest{
    1:required string user_id (api.query="user_id"),
    2:required i64 page_num (api.query="page_num"),
    3:required i64 page_size (api.query="page_size"),
}
struct UploadListResponse{
    1:model.BaseResp base,
    2:model.VideoList data,
}

//热门排行榜
struct RankRequest{
    1:required i64 page_num (api.query="page_num"),
    2:required i64 page_size (api.query="page_size"),
}
struct RankResponse{
    1:model.BaseResp base,
    2:model.VideoList data,
}
//搜索视频
struct QueryRequest{
    1:optional string keywords (api.form="keywords"),
    2:required i64 page_size (api.form="page_size"),
    3:required i64 page_num (api.form="page_num"),
    4:optional i64 from_date (api.form="from_date"),
    5:optional i64 to_date (api.form="to_date"),
    6:optional string username (api.form="username"),
}
struct QueryResponse{
    1:model.BaseResp base,
    2:model.VideoList data,
}

service VideoService{
    FeedResponse Feed(1:FeedRequest req)(api.get="/video/feed"),
    UploadResponse Upload(1:UploadRequest req)(api.post="/video/publish"),
    UploadListResponse UploadList(1:UploadListRequest req)(api.get="/video/list"),
    RankResponse Rank(1:RankRequest req)(api.get="/video/popular"),
    QueryResponse Query(1:QueryRequest req)(api.post="/video/search")
}