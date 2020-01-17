<template>
  <div class="loginContainer">
    <div>
      <div style="font-size:50%;margin-top: 50%;margin-left: 3%">
        <effect-input v-model="value" type="hoshi" label="请输入频道:"></effect-input><span style="margin-left: 4%"></span>
        <button @click="insert()"  v-bind:class="my_cls">加入</button>
      </div>
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
  import 'effect-input/dist/index.css'


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
        my_cls:'btn',
        value:"",
        aletMsg: '', // 弹出框中的提示语
        displayStsates: 'none',


      }
    },
    mounted: function () {

    },
    destroyed: function() {
      //页面销毁时关闭长连接
      //this.websock.close
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
      alertDia (msg) {
        this.displayStsates = 'block'
        this.aletMsg = msg
        // 延迟2秒后消失 自己可以更改时间
        window.setTimeout(() => {
          this.displayStsates = 'none'
        }, 2000)
      },
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
      async insert(){
        var thisVue = this
        if (this.value == ""){
          this.alertDia("频道不能为空")
          return
        }
        var password=prompt("请输入频道密码","123456")
        if (password==null){
          return
        }
        if (password == ""){
          this.alertDia("密码不能为空")
          return
        }

        const data = { channel: this.value,password:password };
        //记得加http
        fetch('http://'+this.url+'/create/channel', {
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
            if (res=="成功"){
              this.$router.push({name:'/login',params:{channel: thisVue.value,password:password}})
            }else{
              this.alertDia(res)
            }
            console.log(res) // res是最终的结果
          })

      },
      closeTip() {
        this.showAlert = false;
      }
    }
  }

</script>

<style lang="scss" scoped>
  @import '../../style/mixin';

  .loginContainer {
    padding-top: 1.95rem;

    p, span, input {
      font-family: Helvetica Neue, Tahoma, Arial;
    }
  }

  .btn {
    width: 100px;
    height: 40px;
    background: gray;
  }

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

</style>
