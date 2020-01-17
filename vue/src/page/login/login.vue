<template>
  <div class="loginContainer">
    <div class="top">
      <p
        style="margin-left: 1%; float:left; width: 20%; line-height: 60px;text-align: center; overflow: hidden;height: 60px;">
        返回(99+)</p>
      <p style="margin-left: 5%;margin-right: 5%; float:left; width: 35%; line-height: 60px; text-align: center">聊天</p>
      <p style="margin-right: 1%; float:right; width:30%;  line-height: 60px;text-align: center" @click="deleteChannel()" >删除频道</p>
      <div style="clear: both"></div>
    </div>

    <div class="mid" id="mid">
      <div class="mid-body">
        <div id="wrapper" style="-webkit-overflow-scrolling : touch; font-size: 50% ">
          <ul style="padding: 0">
            <li v-for="(data,key) in messages">
              <br/> <span style=" word-wrap: break-word;word-break: break-all;overflow: hidden;"><span style="color: red;">消息:</span>{{data}}</span><br/>
              <div style="clear: both">
              </div>
            </li>
            <br/><br/>
          </ul>
        </div>
      </div>
    </div>
    <div class="foot">
      <input class="foot-body"  @keyup.enter="insert()" v-model="message">
      <!--<input type="text" name="inputText" class="inputText" value="1">  <input/>-->
      <input type="button" class="inputSend"  @click="insert()" value="发送">
      <input type="button" class="inputSend"  @click="deleteMeeage()"  value="删除全部">
    </div>

    <!--    提示  -->
    <div class="alsrtInfo" :style="{display: displayStsates}" ref="alertMsg">
      <div class="profPrompt_test">{{aletMsg}}</div>
    </div>

  </div>

</template>

<script >
  import headTop from '../../components/header/head'
  import alertTip from '../../components/common/alertTip'
  import {localapi, proapi, imgBaseUrl} from 'src/config/env'
  import {mapState, mapMutations} from 'vuex'
  import {mobileCode, checkExsis, sendLogin, getcaptchas, accountLogin} from '../../service/getData'

  //window
  /*window.onbeforeunload=function(e){
    var e = window.event||e;
    e.returnValue=("确定离开当前页面吗？");
  }*/

  export default {
    data() {
      return {
        messages: ["开始聊天"],
        url: "47.102.102.47:4000",
        loginWay: false, //登录方式，默认短信登录
        showPassword: false, // 是否显示密码
        phoneNumber: null, //电话号码
        mobileCode: null, //短信验证码
        validate_token: null, //获取短信时返回的验证值，登录时需要
        computedTime: 0, //倒数记时
        userInfo: null, //获取到的用户信息
        userAccount: null, //用户名
        passWord: null, //密码
        captchaCodeImg: null, //验证码地址
        codeNumber: null, //验证码
        showAlert: false, //显示提示组件
        alertText: null, //提示的内容
        websock: null,
        message:"",
        channel:this.$route.params.channel,
        password:this.$route.params.password,
        aletMsg: '', // 弹出框中的提示语
        displayStsates: 'none',

      }
    },
    mounted: function () {
      // alert("进入 websocket")
      // var init = {
      //   body: JSON.stringify({}), // must match 'Content-Type' header
      //   cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
      //   headers: {
      //     'content-type': 'text/plain'
      //   },
      //   method: 'POST', // *GET, POST, PUT, DELETE, etc.
      //   mode: 'cors', // no-cors, cors, *same-origin
      // }
      // fetch(
      //   this.data.url+'/web/socket',
      //   init
      // )
      //   .then(res => alert("请求成功") /*res.json()*///)
      //
      //     .then(data => {
      //       alert("请求成功")
      //     })
      //     .catch(e => console.log('错误:', e))
      //
      //
      //
      // //用户父子组件传值
      // // super(props)
      //
      // this.state = {
      //   title: "标签"
      // }
      // alert("ss")
    },
    destroyed: function() {
      //页面销毁时关闭长连接
      //this.websock.close
    },
    beforeMount(){
      var thisVue = this
      //console.log(this.messages)
      this.websock = new WebSocket('ws://' + this.url + '/ws?channel='+this.channel);
      this.websock.onmessage = function (e) {
        thisVue.messages.push(e.data)
        //alert(thisVue.messages.toString())
      };

      window.setInterval(() => {
        setTimeout(function(){
          thisVue.websock.send("1")
        }, 0)
      }, 1000)

      this.GetMessage()
      //  ws.close()
    },
    created() {


    },
    components: {
      headTop,
      alertTip,
    },
    computed: {
      //判断手机号码
      rightPhoneNumber: function () {
        return /^1\d{10}$/gi.test(this.phoneNumber)
      }
    },
    methods: {
      ...mapMutations([
        'RECORD_USERINFO',
      ]),
      //改变登录方式
      changeLoginWay() {
        this.loginWay = !this.loginWay;
      },
      //是否显示密码
      changePassWordType() {
        this.showPassword = !this.showPassword;
      },
      //获取验证吗，线上环境使用固定的图片，生产环境使用真实的验证码
      async getCaptchaCode() {
        let res = await getcaptchas();
        this.captchaCodeImg = res.code;
      },
      //新增
      async insert(){
      /*  let res = await InsertMessage(this.message);
        if (res.message) {
          alert(res.message)
          return
        }*/
      if (this.message == "") {
        alert("消息不能为空")
        return
      }
      const data = { message: this.message,channel: this.channel ,password:this.password};
      //记得加http
      fetch('http://'+this.url+'/insert/message', {
        method: 'POST',
        headers: {
          'Origin': "HeaderAccessControlAllowOrigin",
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data), // "{"name":"Billy","age":18}"
      }).then((res)=>{
        return res.text() // res.text()是一个Promise对象
      })
        .then((res)=>{
          if (res != "成功"){
            this.alertDia("发送消息失败:"+res)
          }
          console.log(res) // res是最终的结果
        })
        //把输入框里的值变为""
        this.message = [""]
      },

      async GetMessage(){


        var thisVue = this
        const data = { channel: this.channel,password:this.password };
        fetch('http://'+this.url+'/get/message', {
          method: 'POST',
          headers: {
            'Origin': "HeaderAccessControlAllowOrigin",
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(data), // "{"name":"Billy","age":18}"
        }).then((res)=>{
          return res.text() // res.text()是一个Promise对象
        })
          .then((res)=>{
            thisVue.messages = JSON.parse(res)
            console.log(res) // res是最终的结果
          })
      },
      async deleteChannel(){

        var password=prompt("请输入频道密码","123456")
        if (password==null){
          return
        }
        if (password == ""){
          this.alertDia("密码不能为空")
          return
        }

        const data = { channel: this.channel ,password:password};
        var thisVue = this
        await fetch('http://'+this.url+'/delete/channel', {
          method: 'POST',
          headers: {
            'Origin': "HeaderAccessControlAllowOrigin",
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(data),// "{"name":"Billy","age":18}"
        }).then((res)=>{
          return res.text() // res.text()是一个Promise对象
        })
          .then((res)=>{
            if (res != "成功"){
              this.alertDia("删除频道失败:"+res)
            }else{
              this.$router.push({path:'/channel'})
            }
            console.log(res) // res是最终的结果
          })
      },
      async deleteMeeage(){
        const data = { channel: this.channel ,password:this.password};
        var thisVue = this
       await fetch('http://'+this.url+'/delete/message', {
          method: 'POST',
          headers: {
            'Origin': "HeaderAccessControlAllowOrigin",
            'Content-Type': 'application/json',
          },
         body: JSON.stringify(data),// "{"name":"Billy","age":18}"
        }).then((res)=>{
          return res.text() // res.text()是一个Promise对象
        })
          .then((res)=>{
            if (res != "成功"){
              this.alertDia("删除消息失败:",res)
            }
            console.log(res) // res是最终的结果
          })
        await  this.GetMessage()
      },

      //获取短信验证码
      async getVerifyCode() {
        if (this.rightPhoneNumber) {
          this.computedTime = 30;
          this.timer = setInterval(() => {
            this.computedTime--;
            if (this.computedTime == 0) {
              clearInterval(this.timer)
            }
          }, 1000)
          //判读用户是否存在
          let exsis = await checkExsis(this.phoneNumber, 'mobile');
          if (exsis.message) {
            this.showAlert = true;
            this.alertText = exsis.message;
            return
          } else if (!exsis.is_exists) {
            this.showAlert = true;
            this.alertText = '您输入的手机号尚未绑定';
            return
          }
          //发送短信验证码
          let res = await mobileCode(this.phoneNumber);
          if (res.message) {
            this.showAlert = true;
            this.alertText = res.message;
            return
          }
          this.validate_token = res.validate_token;
        }
      },
      alertDia (msg) {
        this.displayStsates = 'block'
        this.aletMsg = msg
        // 延迟2秒后消失 自己可以更改时间
        window.setTimeout(() => {
          this.displayStsates = 'none'
        }, 2000)
      },
      //发送登录信息
      async mobileLogin() {
        if (this.loginWay) {
          if (!this.rightPhoneNumber) {
            this.showAlert = true;
            this.alertText = '手机号码不正确';
            return
          } else if (!(/^\d{6}$/gi.test(this.mobileCode))) {
            this.showAlert = true;
            this.alertText = '短信验证码不正确';
            return
          }
          //手机号登录
          this.userInfo = await sendLogin(this.mobileCode, this.phoneNumber, this.validate_token);
        } else {
          if (!this.userAccount) {
            this.showAlert = true;
            this.alertText = '请输入手机号/邮箱/用户名';
            return
          } else if (!this.passWord) {
            this.showAlert = true;
            this.alertText = '请输入密码';
            return
          } else if (!this.codeNumber) {
            this.showAlert = true;
            this.alertText = '请输入验证码';
            return
          }
          //用户名登录
          this.userInfo = await accountLogin(this.userAccount, this.passWord, this.codeNumber);
        }
        //如果返回的值不正确，则弹出提示框，返回的值正确则返回上一页
        if (!this.userInfo.user_id) {
          this.showAlert = true;
          this.alertText = this.userInfo.message;
          if (!this.loginWay) this.getCaptchaCode();
        } else {
          this.RECORD_USERINFO(this.userInfo);
          this.$router.go(-1);
        }
      },
      closeTip() {
        this.showAlert = false;
      }
    }
  }

</script>

<style lang="scss" scoped>
  @import '../../style/mixin';


  .alsrtInfo{
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 10;
    background: rgba(0, 0, 0, 0.1);
    .profPrompt_test{
      padding: 30px 10px;
      width: 480px;
      overflow: hidden;
      line-height: 28px;
      border: 1px solid #4eb6d3;
      color: #4eb6d3;
      position: absolute;
      background-color: #fbfbfb;
      top: 356px;
      left: 50%;
      font-size: 14px;
      font-family: Gotham-Book;
      opacity: 1;
      /* z-index: 1; */
      text-align: center;
      transform: translate(-50%, -50%);
    }
  }

  .loginContainer {
    padding-top: 1.95rem;

    p, span, input {
      font-family: Helvetica Neue, Tahoma, Arial;
    }
  }

  .change_login {
    position: absolute;
    @include ct;
    right: 0.75rem;
    @include sc(.7rem, #fff);
  }

  .loginForm {
    background-color: #fff;
    margin-top: .6rem;

    .input_container {
      display: flex;
      justify-content: space-between;
      padding: .6rem .8rem;
      border-bottom: 1px solid #f1f1f1;

      input {
        @include sc(.7rem, #666);
      }

      button {
        @include sc(.65rem, #fff);
        font-family: Helvetica Neue, Tahoma, Arial;
        padding: .28rem .4rem;
        border: 1px;
        border-radius: 0.15rem;
      }

      .right_phone_number {
        background-color: #4cd964;
      }
    }

    .phone_number {
      padding: .3rem .8rem;
    }

    .captcha_code_container {
      height: 2.2rem;

      .img_change_img {
        display: flex;
        align-items: center;

        img {
          @include wh(3.5rem, 1.5rem);
          margin-right: .2rem;
        }

        .change_img {
          display: flex;
          flex-direction: 'column';
          flex-wrap: wrap;
          width: 2rem;
          justify-content: center;

          p {
            @include sc(.55rem, #666);
          }

          p:nth-of-type(2) {
            color: #3b95e9;
            margin-top: .2rem;
          }
        }
      }
    }
  }

  .login_tips {
    @include sc(.5rem, red);
    padding: .4rem .6rem;
    line-height: .5rem;

    a {
      color: #3b95e9;
    }
  }

  .login_container {
    margin: 0 .5rem 1rem;
    @include sc(.7rem, #fff);
    background-color: #4cd964;
    padding: .5rem 0;
    border: 1px;
    border-radius: 0.15rem;
    text-align: center;
  }

  .button_switch {
    background-color: #ccc;
    display: flex;
    justify-content: center;
    @include wh(2rem, .7rem);
    padding: 0 .2rem;
    border: 1px;
    border-radius: 0.5rem;
    position: relative;

    .circle_button {
      transition: all .3s;
      position: absolute;
      top: -0.2rem;
      z-index: 1;
      left: -0.3rem;
      @include wh(1.2rem, 1.2rem);
      box-shadow: 0 0.026667rem 0.053333rem 0 rgba(0, 0, 0, .1);
      background-color: #f1f1f1;
      border-radius: 50%;
    }

    .trans_to_right {
      transform: translateX(1.3rem);
    }

    span {
      @include sc(.45rem, #fff);
      transform: translateY(.05rem);
      line-height: .6rem;
    }

    span:nth-of-type(2) {
      transform: translateY(-.08rem);
    }
  }

  .change_to_text {
    background-color: #4cd964;
  }

  .to_forget {
    float: right;
    @include sc(.6rem, #3b95e9);
    margin-right: .3rem;
  }


  * {
    margin: 0;
    border: 0;
  }

  ul, li {
    list-style-type: none;
  }

  body {
    font-family: "Microsoft YaHei";
    max-height: 100%;
    background-color: #EBEBEB;
    overflow: hidden;
    color: #000000;
  }

  .top {
    width: 100%;
    height: 8%;
    z-index: 1000;
    text-align: center;
    color: #ffffff
  }

  .mid {
    max-height: calc(100vh - 120px);
    overflow-x: hidden;
    overflow-y: scroll;
  }

  .mid-body {
    padding: 0 10px;
    padding-bottom: 30px;
  }

  .vmargin {
    margin-bottom: 28vh;
  }

  .chat-left {
    width: 100%;
    min-height: 50px;
    padding: 10px 0;
    text-align: left;
  }

  .chat-left .userimg {
    width: 50px;
    float: left;
    height: 50px;
    border-radius: 50%;
  }

  .chat-left .msg-word {
    margin-left: 5%;
    padding: 15px;
    float: left;
    max-width: 60%;
    background-color: #ffffff;
    border-radius: 10px;
    position: relative;
  }

  .chat-left .msg-word:before {
    content: "";
    position: absolute;
    right: 100%;
    top: 10px;
    width: 0;
    height: 0;
    border-top: 6px solid transparent;
    border-right: 12px solid #ffffff;
    border-bottom: 6px solid transparent;
  }


  .chat-mid {
    padding: 5px;
    width: 100%;
    height: 30px;
    line-height: 30px;
  }

  .chat-mid .time {
    font-size: 10px;
    color: #999;
    text-align: center;
  }


  .chat-right {
    text-align: right;
    width: 100%;
    min-height: 50px;
    padding: 10px 0;
  }

  .chat-right .userimg {
    width: 50px;
    float: right;
    height: 50px;
    border-radius: 50%;
  }

  .chat-right .msg-word {
    margin-right: 5%;
    padding: 15px;
    float: right;
    max-width: 60%;
    border-radius: 10px;
    background-color: #41B883;
    position: relative;
  }

  .showBigImg {

    width: 100vw;
    height: 100vh;
    position: absolute;
    top: 0;
    left: 0;
    filter: alpha(opacity=100);
    opacity: 1;
    -moz-opacity: 1;
    background-color: #222;
    z-index: 9999;
  }

  .bigImgDiv {
    width: 100vw;
    height: 100vh;
    z-index: 10000;
    text-align: center;
    overflow: hidden;
    font-size: 0;
  }

  .bigImgDiv::after {
    content: "";
    height: 100%;
    display: inline-block;
    vertical-align: middle;
  }

  .bigImg {
    width: 100vw;
    height: auto;
    z-index: 10000;
    vertical-align: middle;
  }

  .chat-right .msg-word:before {
    content: "";
    position: absolute;
    left: 100%;
    top: 10px;
    width: 0;
    height: 0;
    border-top: 6px solid transparent;
    border-left: 12px solid #41B883;
    border-bottom: 6px solid transparent;
  }

  .imgMin {
    margin-right: 5%;
    padding: 15px;
    max-width: 60%;
    border-radius: 10px;
    position: relative;
  }

  .imgMax {
    width: 100vw;
    max-width: 100vw;
    height: auto;
  }

  .foot {
    width: 100%;
    /*height: 8%;*/
    position: fixed;
    z-index: 1000;
    background-color: #F5F5F7;
    bottom: 0;
    left: 0;
    text-align: center;
    color: #ffffff
  }

  .foot .foot-body {
    text-align: left;
    padding: 1% 4%;
  }

  .foot .inputText {
    width: 60%;
    height: calc(6vh);
    padding-left: 2px;
  }

  .foot .inputSend {
    margin-left: calc(2vw);
    width: 15%;
    height: calc(6vh);
  }

  .chooseImg {
    padding: 2vh 5vw;
    text-align: left;
    width: 100vw;
    height: 28vh;
    border-top: 1px solid #dddddd;
  }

  .upfilePhoto {
    position: relative;
    display: inline-block;
    background: #D0EEFF;
    border: 1px solid #99D3F5;
    border-radius: 4px;
    padding: 4px 12px;
    overflow: hidden;
    color: #1E88C7;
    text-decoration: none;
    text-indent: 0;
    line-height: 8vh;
  }

  .upfilePhoto input {
    position: absolute;
    font-size: 60px;
    right: 0;
    top: 0;
    opacity: 0;
  }

  .upfilePhoto:hover {
    background: #AADFFD;
    border-color: #78C3F3;
    color: #004974;
    text-decoration: none;
  }

</style>
