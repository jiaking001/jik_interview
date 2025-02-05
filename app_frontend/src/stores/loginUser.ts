import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RootState } from "@/stores/index";
import ACCESS_ENUM from "@/access/accessEnum";
import {DEFAULT_USER} from "@/constants/user";


/**
 * 登录用户全局状态
 */
export const loginUserSlice = createSlice({
    name: "loginUser",
    initialState: DEFAULT_USER,
    reducers: {
        setLoginUser: (state, action: PayloadAction<API.LoginUserVO>) => {
            return {
                ...action.payload,
            };
        },
    },
});

// 修改状态
export const { setLoginUser } = loginUserSlice.actions;

export default loginUserSlice.reducer;
