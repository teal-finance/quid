import{d as N,i as z,j as w,k as F,l as Y,o as i,c as g,a as l,m as Z,v as j,u as r,p as E,q as J,x as D,y as O,b as c,S as A,w as h,z as Q,A as T,B as y,t as B,F as K,e as W,f as x,D as X,N as $,_ as P,E as ee,G as te,H as M,h as L,J as ae,I as H}from"./index-6b0967df.js";import{s as ne,a as C,A as V}from"./ActionButton-36ce7124.js";import{_ as le}from"./SimpleBadge-0cbc5bef.js";import{_ as se}from"./save-outlined-ecd7da7b.js";function ie(){const s=function(){return((1+Math.random())*65536|0).toString(16).substring(1)};return s()+s()+"-"+s()+"-"+s()+"-"+s()+"-"+s()+s()+s()}const oe={class:"inline-block sw-switch",style:{"background-color":"transparent",color:"inherit"}},de=["for"],re={class:"relative"},ce=["id","checked"],G=N({props:{big:{type:Boolean,default:!1},value:{type:Boolean,required:!0}},emits:["update:value"],setup(s,{emit:d}){const m=s,{value:e}=z(m),o=w(!1),f=ie();function u(){d("update:value",o.value)}return F(()=>{o.value=e.value}),Y(()=>m.value,(v,n)=>{o.value=v}),(v,n)=>(i(),g("div",oe,[l("label",{for:"toggle"+r(f),class:"flex items-center cursor-pointer"},[l("div",re,[Z(l("input",{id:"toggle"+r(f),type:"checkbox",class:"sr-only",checked:o.value,"onUpdate:modelValue":n[0]||(n[0]=a=>o.value=a),onChange:n[1]||(n[1]=a=>u())},null,40,ce),[[j,o.value]]),l("div",{class:E(["block rounded-full bg",{big:s.big===!0}])},null,2),l("div",{class:E(["absolute transition rounded-full dot left-1 top-1",{big:s.big===!0}])},null,2)]),J(v.$slots,"default")],8,de)]))}}),ue={class:"flex flex-col space-y-5"},ve=y("\xA0Enable public endpoint"),me={class:"flex flex-row"},pe=["disabled"],_e=N({__name:"AddNamespace",emits:["end"],setup(s,{emit:d}){const m=w(!1),e=D({name:{val:"",isValid:null,validator:v=>v.length>=3},accessTokenTtl:{val:"20m",isValid:!0,validator:v=>v.length>=2},refreshTokenTtl:{val:"24h",isValid:!0,validator:v=>v.length>=2}}),o=O(()=>e.name.isValid===!0&&e.accessTokenTtl.isValid===!0&&e.refreshTokenTtl.isValid===!0);function f(){d("end")}async function u(){try{await Q.post("/admin/namespaces/add",{name:e.name.val,max_ttl:e.accessTokenTtl.val,refresh_max_ttl:e.refreshTokenTtl.val,enable_endpoint:m.value}),d("end"),T.done("Namespace added")}catch(v){console.log(v),T.error("Error creating namespace")}}return(v,n)=>(i(),g("div",ue,[c(r(A),{value:e.name.val,"onUpdate:value":n[0]||(n[0]=a=>e.name.val=a),isvalid:e.name.isValid,"onUpdate:isvalid":n[1]||(n[1]=a=>e.name.isValid=a),validator:e.name.validator,"inline-label":"Name",required:"",autofocus:""},null,8,["value","isvalid","validator"]),c(r(A),{value:e.accessTokenTtl.val,"onUpdate:value":n[2]||(n[2]=a=>e.accessTokenTtl.val=a),isvalid:e.accessTokenTtl.isValid,"onUpdate:isvalid":n[3]||(n[3]=a=>e.accessTokenTtl.isValid=a),validator:e.accessTokenTtl.validator,"inline-label":"Access tokens max time to live",required:""},null,8,["value","isvalid","validator"]),c(r(A),{value:e.refreshTokenTtl.val,"onUpdate:value":n[4]||(n[4]=a=>e.refreshTokenTtl.val=a),isvalid:e.refreshTokenTtl.isValid,"onUpdate:isvalid":n[5]||(n[5]=a=>e.refreshTokenTtl.isValid=a),validator:e.refreshTokenTtl.validator,"inline-label":"Refresh tokens max time to live",required:""},null,8,["value","isvalid","validator"]),c(r(G),{id:"ns-switch",class:"switch-success",value:m.value,"onUpdate:value":n[6]||(n[6]=a=>m.value=a)},{default:h(()=>[ve]),_:1},8,["value"]),l("div",me,[l("button",{class:"w-20 mr-3 btn success",disabled:!r(o),onClick:n[7]||(n[7]=a=>u())},"Save",8,pe),l("button",{class:"w-20 btn warning",onClick:n[8]||(n[8]=a=>f())},"Cancel")])]))}}),fe={class:"flex flex-col p-3 ml-5"},ke=l("span",{class:"pr-2 font-bold"},"Users:",-1),he={class:"mt-2 space-x-2 font-bold"},ge=l("span",{class:"pr-2"},"Groups:",-1),xe=N({__name:"NamespaceInfo",props:{numUsers:{type:Number,required:!0},groups:{type:Array,required:!0}},setup(s){return(d,m)=>(i(),g("div",fe,[l("div",null,[ke,y(" "+B(s.numUsers),1)]),l("div",he,[ge,(i(!0),g(K,null,W(s.groups,e=>(i(),x(le,{key:e.id,class:"inline-block",text:e.name,color:"secondary"},null,8,["text"]))),128))])]))}}),be={class:"inline-block align-middle",preserveAspectRatio:"xMidYMid meet",viewBox:"0 0 24 24",width:"1.2em",height:"1.2em"},ye=X('<g fill="none" stroke="currentColor" stroke-linecap="round" stroke-width="2"><path stroke-dasharray="60" stroke-dashoffset="60" d="M5.63604 5.63603C9.15076 2.12131 14.8492 2.12131 18.364 5.63603C21.8787 9.15075 21.8787 14.8492 18.364 18.364C14.8492 21.8787 9.15076 21.8787 5.63604 18.364C2.12132 14.8492 2.12132 9.15075 5.63604 5.63603Z"><animate fill="freeze" attributeName="stroke-dashoffset" dur="0.5s" values="60;0"></animate></path><path stroke-dasharray="18" stroke-dashoffset="18" d="M6 6L18 18"><animate fill="freeze" attributeName="stroke-dashoffset" begin="0.6s" dur="0.2s" values="18;0"></animate></path></g>',1),Te=[ye];function $e(s,d){return i(),g("svg",be,Te)}var we={name:"line-md-cancel",render:$e};const Ce={class:"inline-block align-middle",preserveAspectRatio:"xMidYMid meet",viewBox:"0 0 24 24",width:"1.2em",height:"1.2em"},Ne=l("path",{fill:"currentColor",d:"M12 20a8 8 0 0 0 8-8a8 8 0 0 0-8-8a8 8 0 0 0-8 8a8 8 0 0 0 8 8m0-18a10 10 0 0 1 10 10a10 10 0 0 1-10 10C6.47 22 2 17.5 2 12A10 10 0 0 1 12 2m.5 5v5.25l4.5 2.67l-.75 1.23L11 13V7h1.5Z"},null,-1),Ve=[Ne];function Ee(s,d){return i(),g("svg",Ce,Ve)}var Me={name:"mdi-clock-outline",render:Ee};const Ae={class:"inline-block align-middle",preserveAspectRatio:"xMidYMid meet",viewBox:"0 0 24 24",width:"1.2em",height:"1.2em"},Ue=l("path",{fill:"currentColor",d:"M21 13.1c-.1 0-.3.1-.4.2l-1 1l2.1 2.1l1-1c.2-.2.2-.6 0-.8l-1.3-1.3c-.1-.1-.2-.2-.4-.2m-1.9 1.8l-6.1 6V23h2.1l6.1-6.1l-2.1-2M12.5 7v5.2l4 2.4l-1 1L11 13V7h1.5M11 21.9c-5.1-.5-9-4.8-9-9.9C2 6.5 6.5 2 12 2c5.3 0 9.6 4.1 10 9.3c-.3-.1-.6-.2-1-.2s-.7.1-1 .2C19.6 7.2 16.2 4 12 4c-4.4 0-8 3.6-8 8c0 4.1 3.1 7.5 7.1 7.9l-.1.2v1.8Z"},null,-1),qe=[Ue];function Re(s,d){return i(),g("svg",Ae,qe)}var Se={name:"mdi-clock-edit-outline",render:Re};const Be={key:0},De=["innerHTML"],Le={key:1},He={class:"ml-2"},I=N({__name:"EditTokenTtl",props:{id:{type:Number,required:!0},ttl:{type:String,required:!0},tokenType:{type:String,required:!0}},emits:["end"],setup(s,{emit:d}){const m=s,e=D({ttlval:{val:m.ttl,isValid:null,validator:n=>n.length>=2}}),o=w(!1);function f(){o.value=!0}function u(){o.value=!1}async function v(){!e.ttlval.isValid||(o.value=!1,m.tokenType=="refresh"?await $.saveMaxRefreshTokenTtl(m.id,e.ttlval.val):await $.saveMaxAccessTokenTtl(m.id,e.ttlval.val),T.done("Ttl changed"),d("end",e.ttlval.val))}return(n,a)=>{const U=Se,q=Me,R=we,S=se;return o.value?(i(),g("div",Le,[c(R,{class:"inline-block txt-neutral",onClick:a[1]||(a[1]=b=>u())}),l("span",He,[c(r(A),{class:"inline-block ttl-inline",value:e.ttlval.val,"onUpdate:value":a[2]||(a[2]=b=>e.ttlval.val=b),isvalid:e.ttlval.isValid,"onUpdate:isvalid":a[3]||(a[3]=b=>e.ttlval.isValid=b),validator:e.ttlval.validator,autofocus:!0,required:""},null,8,["value","isvalid","validator"]),c(S,{class:E(["ml-2 text-xl",e.ttlval.isValid?["txt-success","cursor-pointer"]:["txt-light"]]),onClick:a[4]||(a[4]=b=>v())},null,8,["class"])])])):(i(),g("div",Be,[l("div",{class:"cursor-pointer group",onClick:a[0]||(a[0]=b=>f())},[c(U,{class:"hidden group-hover:inline-block txt-neutral"}),c(q,{class:"group-hover:hidden txt-light"}),l("span",{innerHTML:s.ttl,class:"ml-2"},null,8,De)])]))}}});const Ie=["innerHTML"],Ke=["innerHTML"],Ge=y("Select"),ze=y("Show info"),Fe=y("Hide info"),Ye=y("Show key"),Ze=y("Delete"),je={class:"flex flex-col"},Je={class:"text-xl"},Oe={class:"mt-3"},Qe={class:"flex flex-row mt-5 space-x-2"},We=["onClick"],Xe=N({__name:"NamespaceDatatable",props:{namespaces:{type:Array,required:!0}},emits:["reload"],setup(s,{emit:d}){const m=s,e=w([]),o=w(0),f=ee(),u=D({numUsers:0,groups:new Array});function v(p){T.confirmDelete(`Delete the ${p.name} namespace?`,()=>{$.delete(p.id).then(()=>{T.done("Namespace deleted"),d("reload")})})}async function n(p){const k=await $.fetchRowInfo(p);u.numUsers=k.numUsers,u.groups=k.groups,e.value=m.namespaces.filter(t=>t.id==p),o.value=p}function a(){e.value=[],o.value=0,u.numUsers=0,u.groups=new Array}function U(){f.removeAllGroups()}function q(p){navigator.clipboard.writeText(p),f.removeAllGroups()}function R(p){console.log("Select",p),L.changeNs(p),T.done(`Namespace ${p.name} selected`)}async function S(p,k){const t=await $.getKey(p);console.log("K",t),f.add({severity:"info",summary:k,detail:t})}async function b(p,k){await $.togglePublicEndpoint(p,k),T.done("Endpoint toggled")}return(p,k)=>(i(),g("div",null,[c(r(ne),{value:s.namespaces,class:"main-table p-datatable-sm",expandedRows:e.value,"onUpdate:expandedRows":k[1]||(k[1]=t=>e.value=t),"data-key":"id"},{expansion:h(()=>[c(xe,{"num-users":u.numUsers,groups:u.groups},null,8,["num-users","groups"])]),default:h(()=>[c(r(C),{field:"id",header:"Id"}),c(r(C),{field:"name",header:"Name"}),c(r(C),{field:"publicEndpointEnabled",header:"Public endpoint"},{body:h(t=>[t.data.name!="quid"?(i(),x(r(G),{key:0,label:"Switch",value:t.data.publicEndpointEnabled,"onUpdate:value":_=>t.data.publicEndpointEnabled=_,class:"table-switch switch-secondary dark:switch-primary",onChange:_=>b(t.data.id,Boolean(_))},null,8,["value","onUpdate:value","onChange"])):M("",!0)]),_:1}),c(r(C),{field:"maxTokenTtl",header:"Access token ttl"},{body:h(t=>[t.data.name!="quid"?(i(),x(I,{key:0,id:t.data.id,ttl:t.data.maxTokenTtl,"token-type":"access",onEnd:_=>t.data.maxTokenTtl=_},null,8,["id","ttl","onEnd"])):(i(),g("span",{key:1,class:"ml-6",innerHTML:t.data.maxTokenTtl},null,8,Ie))]),_:1}),c(r(C),{field:"maxRefreshTokenTtl",header:"Refresh token ttl"},{body:h(t=>[t.data.name!="quid"?(i(),x(I,{key:0,id:t.data.id,ttl:t.data.maxRefreshTokenTtl,"token-type":"refresh",onEnd:_=>t.data.maxRefreshTokenTtl=_},null,8,["id","ttl","onEnd"])):(i(),g("span",{key:1,class:"ml-6",innerHTML:t.data.maxRefreshTokenTtl},null,8,Ke))]),_:1}),c(r(C),{field:"actions"},{body:h(t=>[t.data.name!="quid"?(i(),x(V,{key:0,onClick:_=>R(t.data),class:E(t.data.name!="quid"?"mr-2":""),disabled:t.data.id==r(L).namespace.value.id},{default:h(()=>[Ge]),_:2},1032,["onClick","class","disabled"])):M("",!0),o.value!=t.data.id?(i(),x(V,{key:1,onClick:_=>n(t.data.id)},{default:h(()=>[ze]),_:2},1032,["onClick"])):(i(),x(V,{key:2,onClick:k[0]||(k[0]=_=>a())},{default:h(()=>[Fe]),_:1})),t.data.name!="quid"?(i(),x(V,{key:3,class:"ml-2",onClick:_=>S(t.data.id,t.data.name)},{default:h(()=>[Ye]),_:2},1032,["onClick"])):M("",!0),t.data.name!="quid"?(i(),x(V,{key:4,type:"delete",class:"ml-2",onClick:_=>v(t.data)},{default:h(()=>[Ze]),_:2},1032,["onClick"])):M("",!0)]),_:1})]),_:1},8,["value","expandedRows"]),c(r(te),{position:"top-center"},{message:h(t=>[l("div",je,[l("div",null,[l("div",Je,"Namespace "+B(t.message.summary),1),l("div",Oe,B(t.message.detail),1)]),l("div",Qe,[l("button",{class:"btn primary",onClick:_=>q(t.message.detail)},"Copy",8,We),l("button",{class:"btn",onClick:k[2]||(k[2]=_=>U())},"Close")])])]),_:1})]))}});var Pe=P(Xe,[["__scopeId","data-v-641730ab"]]);const et={class:"text-3xl txt-primary dark:txt-light"},tt=y(" Namespaces "),at={class:"p-5 mt-3 border bord-lighter w-96"},nt=l("div",{class:"text-xl"},"Add a namespace",-1),dt=N({__name:"NamespaceView",setup(s){const d=w(!0),m=w(new Array);function e(){d.value=!d.value,o()}async function o(){const f=await $.fetchAll();m.value=Array.from(f)}return ae(()=>o()),(f,u)=>(i(),g(K,null,[l("div",et,[tt,l("button",{class:"ml-3 text-2xl border-none btn focus:outline-none txt-neutral",onClick:u[0]||(u[0]=v=>d.value=!d.value)},[d.value===!0?(i(),x(r(H),{key:0,icon:"fa6-solid:plus"})):(i(),x(r(H),{key:1,icon:"fa6-solid:minus"}))])]),l("div",{class:E([{"slide-y":!0,slideup:d.value===!0,slidedown:d.value===!1},"mb-8"])},[l("div",at,[nt,c(_e,{class:"mt-5",onEnd:u[1]||(u[1]=v=>e())})])],2),c(Pe,{namespaces:m.value,onReload:u[2]||(u[2]=v=>o())},null,8,["namespaces"])],64))}});export{dt as default};
