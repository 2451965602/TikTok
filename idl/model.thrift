namespace go model

struct User {
    1:required i64 id,
    2:required string username,
    3:optional string password,
    4:required string avatar_url,
    5:optional string opt_secret,
    6:optional string created_at,
    7:optional string updated_at,
    8:optional string deleted_at,
}

struct UserInfo {
    1:required string id,
    2:required string username,
    3:required string avatar_url,
    4:optional string created_at,
    5:optional string updated_at,
    6:optional string deleted_at,
}

struct Video {
    1:required i64 id,
    2:required string user_id,
    3:required string video_url,
    4:required string cover_url,
    5:required string title,
    6:required string description,
    7:required i64 visit_count,
    8:required i64 like_count,
    9:required i64 comment_count,
    10:required string created_at,
    11:required string updated_at,
    12:required string deleted_at,
}

struct Comment {
    1:required i64 id,
    2:required string user_id,
    3:required string root_id,
    4:required string video_id,
    5:required string content,
    6:required string created_at,
    7:required string updated_at,
    8:required string deleted_at,
}

struct like {
    1:required string user_id,
    2:required string root_id,
    3:required string video_id,
}

struct UserList{
    1:required list<User> items,
    2:required i64 total,
}

struct UserInfoList{
    1:required list<UserInfo> items,
    2:required i64 total,
}

struct VideoList{
    1:required list<Video> items,
    2:required i64 total,
}

struct CommentList{
    1:required list<Comment> items,
    2:required i64 total,
}

struct LikeList{
    1:required list<Video> items,
    2:required i64 total,
}

struct MFA{
    1:required string secret,
    2:required string qrcode,
}

struct Social{
    1:required string user_id,
    2:required string to_user_id,
    3:required i64 status,
}

struct BaseResp{
    1:required i64 code,
    2:required string msg,
}