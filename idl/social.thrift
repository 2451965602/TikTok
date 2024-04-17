namespace go social

include "model.thrift"

//关注操作
struct StarRequest{
    1:required string to_user_id (api.form="to_user_id"),
    2:required i64 action_type (api.form="action_type"),
}

struct StarResponse{
    1:model.BaseResp base,
}

//关注列表
struct StarListRequest{
    1:required string user_id (api.query="user_id"),
    2:required i64 page_size (api.query="page_size"),
    3:required i64 page_num (api.query="page_num"),
}

struct StarListResponse{
    1:model.BaseResp base,
    2:model.UserInfoList data,
}

//粉丝列表
struct FanListRequest{
    1:required string user_id (api.query="user_id"),
    2:required i64 page_size (api.query="page_size"),
    3:required i64 page_num (api.query="page_num"),
}

struct FanListResponse{
    1:model.BaseResp base,
    2:model.UserInfoList data,
}

//好友列表
struct FriendListRequest{
    1:required i64 page_size (api.query="page_size"),
    2:required i64 page_num (api.query="page_num"),
}

struct FriendListResponse{
    1:model.BaseResp base,
    2:model.UserInfoList data,
}


service socialService{
    StarResponse Star(1:StarRequest req)(api.post="/relation/action"),
    StarListResponse StarList(1:StarListRequest req)(api.get="/following/list"),
    FanListResponse FanList(1:FanListRequest req)(api.get="/follower/list"),
    FriendListResponse FriendList(1:FriendListRequest req)(api.get="/friends/list")
}