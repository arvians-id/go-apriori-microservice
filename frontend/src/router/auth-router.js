import Login from "@/views/jwt/Login";
import Register from "@/views/jwt/Register";
import ForgotPassword from "@/views/jwt/ForgotPassword";
import ResetPassword from "@/views/jwt/ResetPassword";

const authRouter = [
    { path: "/jwt/login", name: "auth.login", component: Login },
    { path: "/jwt/register", name: "auth.register", component: Register },
    { path: "/jwt/forgot-password", name: "auth.forgot-password", component: ForgotPassword },
    { path: "/jwt/reset-password", name: "auth.reset-password", component: ResetPassword },
]

export default authRouter;