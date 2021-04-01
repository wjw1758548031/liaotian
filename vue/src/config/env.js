/**
 * 配置编译环境和线上环境之间的切换
 * 
 * baseUrl: 域名地址
 * routerMode: 路由模式
 * imgBaseUrl: 图片所在域名地址
 * 
 */

let baseUrl = ''; 
let routerMode = 'hash';
let imgBaseUrl = '';


if (process.env.NODE_ENV == 'development') {
    imgBaseUrl = '/img/';

}else if(process.env.NODE_ENV == 'production'){
	  baseUrl = '//127.0.0.1:4000';
    imgBaseUrl = '//elm.cangdu.org/img/';
}

export {
	baseUrl,
	routerMode,
	imgBaseUrl,
}