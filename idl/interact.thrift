namespace go interact

include "model.thrift"

//点赞操作
struct LikeRequest{
    1:optional string video_id,
    2:optional string comment_id,
    3:required string action_type,
}

struct LikeResponse{
    1:model.BaseResp base,
}

//点赞列表
struct LikeListRequest{
    1:required string user_id,
    2:required i64 page_size,
    3:required i64 page_num,
}

struct LikeListResponse{
    1:model.BaseResp base,
    2:model.LikeList data,
}

//评论
struct CommentRequest{
    1:optional string video_id,
    2:optional string comment_id,
    3:required string content,
}

struct CommentResponse{
    1:model.BaseResp base,
}

//评论列表
struct CommentListRequest{
    1:required string video_id,
    2:required i64 page_size,
    3:required i64 page_num,
}

struct CommentListResponse{
    1:model.BaseResp base,
    2:model.CommentList data,
}

//删除评论
struct DeleteCommentRequest{
    1:required string comment_id,

}

struct DeleteCommentResponse{
    1:model.BaseResp base,
}

service interactService{
    LikeResponse Like(1:LikeRequest req)(api.post="/like/action"),
    LikeListResponse LikeList(1:LikeListRequest req)(api.get="/like/list"),
    CommentResponse Comment(1:CommentRequest req)(api.post="/comment/publish"),
    CommentListResponse CommentList(1:CommentListRequest req)(api.get="/comment/list")
    DeleteCommentResponse DeleteComment(1:DeleteCommentRequest req)(api.delete="/comment/delete")
}