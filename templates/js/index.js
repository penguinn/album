var domain = '';
var app = new Vue({
    el: '#app',
    authCodeImg: '',
    data: {
        message: 'Hello Album',
        registerInfo: {
            loginName: '',
            password: '',
            confirm:'',
        },
        loginInfo: {
            loginName: '',
            password: ''
        },
    },
    mounted: function () {
    },
    methods: {
        errorNotify(title,message) {
            this.$notify.error({
                title: title,
                message: message
            });
        },
  
        successNotify(title,message) {
          this.$notify({
            title: title,
            message: message,
            type: 'success'
          });
        },

        warnMessage(message) {
            this.$message({
                message: message,
                type: 'warning'
            });
        },
        
        switchLogin:function(){
            document.getElementById("registerForm").style.display = 'none';
            $("#loginForm").show();
        },

        switchRegister:function(){
            document.getElementById("loginForm").style.display = 'none';
            this.getAuthCode(1);
            $("#registerForm").show();
        },

        getAuthCode:function(type){
            let vm = this;
            $.ajax({
                type: 'POST',
                url: domain + '/api/check/image-auth',
                dataType: "text",
                data: {
                    type: type,
                },
                success: function (res) {
                    vm.authCodeImg = res;
                    if(vm.authCodeImg) {
                        document.getElementById('authImage').src.dataType = 'image/png';
                        document.getElementById('authImage').src.data = vm.authCodeImg;
                    }else {
                        document.getElementById('image').style.display = "none";
                    }
                },
                error:function (res) {
                    vm.errorNotify("验证码获取失败2",'');
                }

            });
        },


        login:function () {
            let vm = this;
            let logName = $.trim(vm.loginInfo.loginName);
            let logPassword = $.trim(vm.loginInfo.password);
            if(logName == null || logName.length == 0){
                vm.warnMessage('抱歉, 请输入用户名~');
            }
            if(logPassword == null || logPassword.length == 0){
                vm.warnMessage('抱歉, 请输入密码~');
            }
            let logPasswordMd5 = md5(logPassword);
            $.ajax({
                type: 'POST',
                url: domain + '/user/login',
                dataType: "json",
                data: {
                    loginName: logName,
                    password: logPasswordMd5,
                },
                success: function (res) {
                    if(res.succ == true){
                        vm.successNotify("登录成功","欢迎你: " + logName);
                    }else{
                        vm.errorNotify("登录失败",res.msg);
                    }
                },
                error:function (res) {
                    vm.errorNotify("登录失败","我们错了: " + logName);
                }

            });
        },

        register:function () {
            let vm = this;
            let regName = $.trim(vm.registerInfo.loginName);
            let regPassword = $.trim(vm.registerInfo.password);
            let regConfirm = $.trim(vm.registerInfo.confirm);
            if(regName == null || regName.length == 0){
                vm.warnMessage('抱歉, 请输入用户名~');
            }
            if(regPassword == null || regPassword.length == 0){
                vm.warnMessage('抱歉, 请输入密码~');
            }

            if(regConfirm == null || regConfirm.length == 0){
                vm.warnMessage('抱歉, 请输入验证码~');
            }

            let regPasswordMd5 = md5(regPassword);

            $.ajax({
                type: 'POST',
                url: domain + '/api/user/sign-in',
                dataType: "json",
                data: {
                    username: regName,
                    password: regPasswordMd5,
                    authCode: regConfirm,
                },
                success: function (res) {
                    if(res.succ == true){
                        vm.successNotify("注册成功","欢迎你: " + regName);
                    }else{
                        vm.errorNotify("注册失败",res.msg);
                    }
                },
                error:function (res) {
                    vm.errorNotify("注册失败","我们错了: " + regName);
                }

            });
        },

      }
})