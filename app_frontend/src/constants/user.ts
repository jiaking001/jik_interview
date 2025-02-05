import ACCESS_ENUM from "@/access/accessEnum";

// 默认用户
export const
    DEFAULT_USER: API.LoginUserVO = {
        UserName: "未登录",
        UserProfile: "暂无简介",
        UserAvatar: "/assets/notLoginUser.png",
        UserRole: ACCESS_ENUM.NOT_LOGIN,
    };