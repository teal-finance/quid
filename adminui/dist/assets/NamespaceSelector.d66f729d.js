import{N as d,u as f,n as h}from"./index.3f7e8321.js";import{_ as x}from"./SimpleBadge.677b0a98.js";import{o as e,b as t,a as r,d as u,f as v,n as y,x as i,L as g,c as m,F as N,e as k,g as A}from"./vendor.ef2cb0e3.js";const C={class:"inline-block align-middle",width:"1.2em",height:"1.2em",preserveAspectRatio:"xMidYMid meet",viewBox:"0 0 24 24"},w=r("path",{fill:"currentColor",d:"M12 2A10 10 0 1 0 22 12A10 10 0 0 0 12 2Zm0 18a8 8 0 1 1 8-8A8 8 0 0 1 12 20Z",opacity:".5"},null,-1),B=r("path",{fill:"currentColor",d:"M20 12h2A10 10 0 0 0 12 2V4A8 8 0 0 1 20 12Z"},[r("animateTransform",{attributeName:"transform",dur:"1s",from:"0 12 12",repeatCount:"indefinite",to:"360 12 12",type:"rotate"})],-1),M=[w,B];function V(n,o){return e(),t("svg",C,M)}var $={name:"eos-icons-loading",render:V};const b=u({props:{small:{type:Boolean,default:!1}},setup(n){return(o,l)=>{const a=$;return e(),t("div",{class:y(["w-full text-center",n.small?"py-8":"pt-24"])},[v(a,{class:"text-5xl txt-light"})],2)}}}),L={key:1,class:"flex flex-wrap space-x-1"},R=u({emits:["selectns"],setup(n,{emit:o}){const l=i(!1),a=i(new Array);async function p(){const s=await d.fetchAll();a.value=Array.from(s)}function _(s){f.changeNs(s),o("selectns"),h.done("Namespace selected")}return g(()=>p()),(s,Z)=>(e(),t("div",null,[l.value?(e(),m(b,{key:0})):(e(),t("div",L,[(e(!0),t(N,null,k(a.value,c=>(e(),t("div",null,[c.name!="quid"?(e(),m(x,{key:0,text:c.name,class:"mr-2 cursor-pointer primary",onClick:F=>_(c)},null,8,["text","onClick"])):A("",!0)]))),256))]))]))}});export{R as _,b as a};
