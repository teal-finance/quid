import{d as A,x as C,l as S,m as G,o as i,b as h,f as o,u as t,k as M,w as k,a as l,q as g,t as L,F as K,e as Y,c as x,n as B,A as Z,g as E,G as z,L as j,I as D}from"./vendor.ef2cb0e3.js";import{_ as F}from"./switch.es.3ae85dac.js";import{r as J,n as b,N as $,_ as O,u as H}from"./index.3f7e8321.js";import{s as w,A as V,a as Q}from"./ActionButton.c40b50fa.js";import{_ as W}from"./SimpleBadge.677b0a98.js";import{_ as X}from"./save-outlined.59bf967f.js";const P={class:"flex flex-col space-y-5"},ee=g("\xA0Enable public endpoint"),te={class:"flex flex-row"},ae=["disabled"],ne=A({emits:["end"],setup(_,{emit:d}){const m=C(!1),e=S({name:{val:"",isValid:null,validator:u=>u.length>=3},accessTokenTtl:{val:"20m",isValid:!0,validator:u=>u.length>=2},refreshTokenTtl:{val:"24h",isValid:!0,validator:u=>u.length>=2}}),v=G(()=>e.name.isValid===!0&&e.accessTokenTtl.isValid===!0&&e.refreshTokenTtl.isValid===!0);function y(){d("end")}async function r(){try{await J.post("/admin/namespaces/add",{name:e.name.val,max_ttl:e.accessTokenTtl.val,refresh_max_ttl:e.refreshTokenTtl.val,enable_endpoint:m.value}),d("end"),b.done("Namespace added")}catch(u){console.log(u),b.error("Error creating namespace")}}return(u,s)=>(i(),h("div",P,[o(t(M),{value:t(e).name.val,"onUpdate:value":s[0]||(s[0]=n=>t(e).name.val=n),isvalid:t(e).name.isValid,"onUpdate:isvalid":s[1]||(s[1]=n=>t(e).name.isValid=n),validator:t(e).name.validator,"inline-label":"Name",required:"",autofocus:""},null,8,["value","isvalid","validator"]),o(t(M),{value:t(e).accessTokenTtl.val,"onUpdate:value":s[2]||(s[2]=n=>t(e).accessTokenTtl.val=n),isvalid:t(e).accessTokenTtl.isValid,"onUpdate:isvalid":s[3]||(s[3]=n=>t(e).accessTokenTtl.isValid=n),validator:t(e).accessTokenTtl.validator,"inline-label":"Access tokens max time to live",required:""},null,8,["value","isvalid","validator"]),o(t(M),{value:t(e).refreshTokenTtl.val,"onUpdate:value":s[4]||(s[4]=n=>t(e).refreshTokenTtl.val=n),isvalid:t(e).refreshTokenTtl.isValid,"onUpdate:isvalid":s[5]||(s[5]=n=>t(e).refreshTokenTtl.isValid=n),validator:t(e).refreshTokenTtl.validator,"inline-label":"Refresh tokens max time to live",required:""},null,8,["value","isvalid","validator"]),o(t(F),{id:"ns-switch",class:"switch-success",value:m.value,"onUpdate:value":s[6]||(s[6]=n=>m.value=n)},{default:k(()=>[ee]),_:1},8,["value"]),l("div",te,[l("button",{class:"w-20 mr-3 btn success",disabled:!t(v),onClick:s[7]||(s[7]=n=>r())},"Save",8,ae),l("button",{class:"w-20 btn warning",onClick:s[8]||(s[8]=n=>y())},"Cancel")])]))}}),le={class:"flex flex-col p-3 ml-5"},se=l("span",{class:"pr-2 font-bold"},"Users:",-1),ie={class:"mt-2 space-x-2 font-bold"},oe=l("span",{class:"pr-2"},"Groups:",-1),de=A({props:{numUsers:{type:Number,required:!0},groups:{type:Array,required:!0}},setup(_){return(d,m)=>(i(),h("div",le,[l("div",null,[se,g(" "+L(_.numUsers),1)]),l("div",ie,[oe,(i(!0),h(K,null,Y(_.groups,e=>(i(),x(W,{key:e.id,class:"inline-block",text:e.name,color:"secondary"},null,8,["text"]))),128))])]))}}),re={width:"1.2em",height:"1.2em",preserveAspectRatio:"xMidYMid meet",viewBox:"0 0 24 24"},ce=l("g",{fill:"none",stroke:"currentColor","stroke-linecap":"round","stroke-width":"2"},[l("path",{d:"M5.63604 5.63603C9.15076 2.12131 14.8492 2.12131 18.364 5.63603C21.8787 9.15075 21.8787 14.8492 18.364 18.364C14.8492 21.8787 9.15076 21.8787 5.63604 18.364C2.12132 14.8492 2.12132 9.15075 5.63604 5.63603Z",class:"il-md-length-70 il-md-duration-4 il-md-delay-0"}),l("path",{d:"M6 6L18 18",class:"il-md-length-25 il-md-duration-2 il-md-delay-5"})],-1),ue=[ce];function me(_,d){return i(),h("svg",re,ue)}var ve={name:"line-md-cancel",render:me};const pe={class:"inline-block align-middle",width:"1.2em",height:"1.2em",preserveAspectRatio:"xMidYMid meet",viewBox:"0 0 24 24"},_e=l("path",{fill:"currentColor",d:"M12 20a8 8 0 0 0 8-8a8 8 0 0 0-8-8a8 8 0 0 0-8 8a8 8 0 0 0 8 8m0-18a10 10 0 0 1 10 10a10 10 0 0 1-10 10C6.47 22 2 17.5 2 12A10 10 0 0 1 12 2m.5 5v5.25l4.5 2.67l-.75 1.23L11 13V7h1.5Z"},null,-1),fe=[_e];function ke(_,d){return i(),h("svg",pe,fe)}var he={name:"mdi-clock-outline",render:ke};const xe={class:"inline-block align-middle",width:"1.2em",height:"1.2em",preserveAspectRatio:"xMidYMid meet",viewBox:"0 0 24 24"},ye=l("path",{fill:"currentColor",d:"M21 13.1c-.1 0-.3.1-.4.2l-1 1l2.1 2.1l1-1c.2-.2.2-.6 0-.8l-1.3-1.3c-.1-.1-.2-.2-.4-.2m-1.9 1.8l-6.1 6V23h2.1l6.1-6.1l-2.1-2M12.5 7v5.2l4 2.4l-1 1L11 13V7h1.5M11 21.9c-5.1-.5-9-4.8-9-9.9C2 6.5 6.5 2 12 2c5.3 0 9.6 4.1 10 9.3c-.3-.1-.6-.2-1-.2s-.7.1-1 .2C19.6 7.2 16.2 4 12 4c-4.4 0-8 3.6-8 8c0 4.1 3.1 7.5 7.1 7.9l-.1.2v1.8Z"},null,-1),Te=[ye];function ge(_,d){return i(),h("svg",xe,Te)}var be={name:"mdi-clock-edit-outline",render:ge};const $e={key:0},we=["innerHTML"],Ce={key:1},Ve={class:"ml-2"},I=A({props:{id:{type:Number,required:!0},ttl:{type:String,required:!0},tokenType:{type:String,required:!0}},emits:["end"],setup(_,{emit:d}){const m=_,e=S({ttlval:{val:m.ttl,isValid:null,validator:s=>s.length>=2}}),v=C(!1);function y(){v.value=!0}function r(){v.value=!1}async function u(){!e.ttlval.isValid||(v.value=!1,m.tokenType=="refresh"?await $.saveMaxRefreshTokenTtl(m.id,e.ttlval.val):await $.saveMaxAccessTokenTtl(m.id,e.ttlval.val),b.done("Ttl changed"),d("end",e.ttlval.val))}return(s,n)=>{const N=be,U=he,q=ve,R=X;return v.value?(i(),h("div",Ce,[o(q,{class:"inline-block txt-neutral",onClick:n[1]||(n[1]=T=>r())}),l("span",Ve,[o(t(M),{class:"inline-block ttl-inline",value:t(e).ttlval.val,"onUpdate:value":n[2]||(n[2]=T=>t(e).ttlval.val=T),isvalid:t(e).ttlval.isValid,"onUpdate:isvalid":n[3]||(n[3]=T=>t(e).ttlval.isValid=T),validator:t(e).ttlval.validator,autofocus:!0,required:""},null,8,["value","isvalid","validator"]),o(R,{class:B(["ml-2 text-xl",t(e).ttlval.isValid?["txt-success","cursor-pointer"]:["txt-light"]]),onClick:n[4]||(n[4]=T=>u())},null,8,["class"])])])):(i(),h("div",$e,[l("div",{class:"cursor-pointer group",onClick:n[0]||(n[0]=T=>y())},[o(N,{class:"hidden group-hover:inline-block txt-neutral"}),o(U,{class:"group-hover:hidden txt-light"}),l("span",{innerHTML:_.ttl,class:"ml-2"},null,8,we)])]))}}});const Ae=["innerHTML"],Ee=["innerHTML"],Me=g("Select"),Ne=g("Show info"),Ue=g("Hide info"),qe=g("Show key"),Re=g("Delete"),Le={class:"flex flex-col"},Se={class:"text-xl"},Be={class:"mt-3"},De={class:"flex flex-row mt-5 space-x-2"},He=["onClick"],Ie=A({props:{namespaces:{type:Array,required:!0}},emits:["reload"],setup(_,{emit:d}){const m=_,e=C([]),v=C(0),y=Z(),r=S({numUsers:0,groups:new Array});function u(c){b.confirmDelete(`Delete the ${c.name} namespace?`,()=>{$.delete(c.id).then(()=>{b.done("Namespace deleted"),d("reload")})})}async function s(c){const f=await $.fetchRowInfo(c);r.numUsers=f.numUsers,r.groups=f.groups,e.value=m.namespaces.filter(a=>a.id==c),v.value=c}function n(){e.value=[],v.value=0,r.numUsers=0,r.groups=new Array}function N(){y.removeAllGroups()}function U(c){navigator.clipboard.writeText(c),y.removeAllGroups()}function q(c){console.log("Select",c),H.changeNs(c),b.done(`Namespace ${c.name} selected`)}async function R(c,f){const a=await $.getKey(c);console.log("K",a),y.add({severity:"info",summary:f,detail:a})}async function T(c,f){await $.togglePublicEndpoint(c,f),b.done("Endpoint toggled")}return(c,f)=>(i(),h("div",null,[o(t(Q),{value:_.namespaces,class:"main-table p-datatable-sm",expandedRows:e.value,"onUpdate:expandedRows":f[1]||(f[1]=a=>e.value=a),"data-key":"id"},{expansion:k(()=>[o(de,{"num-users":t(r).numUsers,groups:t(r).groups},null,8,["num-users","groups"])]),default:k(()=>[o(t(w),{field:"id",header:"Id"}),o(t(w),{field:"name",header:"Name"}),o(t(w),{field:"publicEndpointEnabled",header:"Public endpoint"},{body:k(a=>[a.data.name!="quid"?(i(),x(t(F),{key:0,label:"Switch",value:a.data.publicEndpointEnabled,"onUpdate:value":p=>a.data.publicEndpointEnabled=p,class:"table-switch switch-secondary dark:switch-primary",onChange:p=>T(a.data.id,Boolean(p))},null,8,["value","onUpdate:value","onChange"])):E("",!0)]),_:1}),o(t(w),{field:"maxTokenTtl",header:"Access token ttl"},{body:k(a=>[a.data.name!="quid"?(i(),x(I,{key:0,id:a.data.id,ttl:a.data.maxTokenTtl,"token-type":"access",onEnd:p=>a.data.maxTokenTtl=p},null,8,["id","ttl","onEnd"])):(i(),h("span",{key:1,class:"ml-6",innerHTML:a.data.maxTokenTtl},null,8,Ae))]),_:1}),o(t(w),{field:"maxRefreshTokenTtl",header:"Refresh token ttl"},{body:k(a=>[a.data.name!="quid"?(i(),x(I,{key:0,id:a.data.id,ttl:a.data.maxRefreshTokenTtl,"token-type":"refresh",onEnd:p=>a.data.maxRefreshTokenTtl=p},null,8,["id","ttl","onEnd"])):(i(),h("span",{key:1,class:"ml-6",innerHTML:a.data.maxRefreshTokenTtl},null,8,Ee))]),_:1}),o(t(w),{field:"actions"},{body:k(a=>[a.data.name!="quid"?(i(),x(V,{key:0,onClick:p=>q(a.data),class:B(a.data.name!="quid"?"mr-2":""),disabled:a.data.id==t(H).namespace.value.id},{default:k(()=>[Me]),_:2},1032,["onClick","class","disabled"])):E("",!0),v.value!=a.data.id?(i(),x(V,{key:1,onClick:p=>s(a.data.id)},{default:k(()=>[Ne]),_:2},1032,["onClick"])):(i(),x(V,{key:2,onClick:f[0]||(f[0]=p=>n())},{default:k(()=>[Ue]),_:1})),a.data.name!="quid"?(i(),x(V,{key:3,class:"ml-2",onClick:p=>R(a.data.id,a.data.name)},{default:k(()=>[qe]),_:2},1032,["onClick"])):E("",!0),a.data.name!="quid"?(i(),x(V,{key:4,type:"delete",class:"ml-2",onClick:p=>u(a.data)},{default:k(()=>[Re]),_:2},1032,["onClick"])):E("",!0)]),_:1})]),_:1},8,["value","expandedRows"]),o(t(z),{position:"top-center"},{message:k(a=>[l("div",Le,[l("div",null,[l("div",Se,"Namespace "+L(a.message.summary),1),l("div",Be,L(a.message.detail),1)]),l("div",De,[l("button",{class:"btn primary",onClick:p=>U(a.message.detail)},"Copy",8,He),l("button",{class:"btn",onClick:f[2]||(f[2]=p=>N())},"Close")])])]),_:1})]))}});var Ke=O(Ie,[["__scopeId","data-v-641730ab"]]);const Fe={class:"text-3xl txt-primary dark:txt-light"},Ge=g(" Namespaces "),Ye={class:"p-5 mt-3 border bord-lighter w-96"},Ze=l("div",{class:"text-xl"},"Add a namespace",-1),Xe=A({setup(_){const d=C(!0),m=C(new Array);function e(){d.value=!d.value,v()}async function v(){const y=await $.fetchAll();m.value=Array.from(y)}return j(()=>v()),(y,r)=>(i(),h(K,null,[l("div",Fe,[Ge,l("button",{class:"ml-3 text-2xl border-none btn focus:outline-none txt-neutral",onClick:r[0]||(r[0]=u=>d.value=!d.value)},[d.value===!0?(i(),x(t(D),{key:0,icon:"fa6-solid:plus"})):(i(),x(t(D),{key:1,icon:"fa6-solid:minus"}))])]),l("div",{class:B([{"slide-y":!0,slideup:d.value===!0,slidedown:d.value===!1},"mb-8"])},[l("div",Ye,[Ze,o(ne,{class:"mt-5",onEnd:r[1]||(r[1]=u=>e())})])],2),o(Ke,{namespaces:m.value,onReload:r[2]||(r[2]=u=>v())},null,8,["namespaces"])],64))}});export{Xe as default};
