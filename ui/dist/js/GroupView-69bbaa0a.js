import{d as w,j as k,o as l,c,b as d,w as x,u as n,B as $,A as f,a5 as C,x as N,y as V,S as A,a as u,z as S,h as i,J as G,f as v,I as g,H as y,p as h,F as B}from"./index-6b0967df.js";import{_ as D}from"./NamespaceSelector-988775c6.js";import{s as R,a as b,A as F}from"./ActionButton-36ce7124.js";import"./SimpleBadge-0cbc5bef.js";const I=$("Delete"),U=w({__name:"GroupDatatable",props:{groups:{type:Array,required:!0}},emits:["reload"],setup(_,{emit:s}){const t=k([]);function m(o){f.confirmDelete(`Delete the ${o.name} group?`,()=>{C.delete(o.id).then(()=>{f.done("Group deleted"),s("reload")})})}return(o,r)=>(l(),c("div",null,[d(n(R),{value:_.groups,class:"main-table p-datatable-sm",expandedRows:t.value,"onUpdate:expandedRows":r[0]||(r[0]=e=>t.value=e),"data-key":"id"},{default:x(()=>[d(n(b),{field:"id",header:"Id"}),d(n(b),{field:"name",header:"Name"}),d(n(b),{field:"actions"},{body:x(e=>[d(F,{type:"delete",class:"ml-2",onClick:a=>m(e.data)},{default:x(()=>[I]),_:2},1032,["onClick"])]),_:1})]),_:1},8,["value","expandedRows"])]))}}),q={class:"flex flex-col space-y-5"},E={class:"flex flex-row"},z=["disabled"],j=w({__name:"AddGroup",emits:["end"],setup(_,{emit:s}){const t=N({name:{val:"",isValid:null,validator:e=>e.length>=3}}),m=V(()=>t.name.isValid===!0);function o(){s("end")}async function r(){try{await S.post(i.adminUrl+"/groups/add",{name:t.name.val,namespace_id:i.namespace.value.id}),s("end"),f.done("Group added")}catch(e){console.log(e),f.error("Error creating group")}}return(e,a)=>(l(),c("div",q,[d(n(A),{value:t.name.val,"onUpdate:value":a[0]||(a[0]=p=>t.name.val=p),isvalid:t.name.isValid,"onUpdate:isvalid":a[1]||(a[1]=p=>t.name.isValid=p),validator:t.name.validator,"inline-label":"Name",required:"",autofocus:""},null,8,["value","isvalid","validator"]),u("div",E,[u("button",{class:"w-20 mr-3 btn success",disabled:!n(m),onClick:a[2]||(a[2]=p=>r())},"Save",8,z),u("button",{class:"w-20 btn warning",onClick:a[3]||(a[3]=p=>o())},"Cancel")])]))}}),H={class:"text-3xl txt-primary dark:txt-light"},J=$(" Groups "),M={class:"p-5 mt-3 border bord-lighter w-96"},T=u("div",{class:"text-xl"},"Add a group",-1),K={key:1,class:"w-full"},L=u("div",{class:"mt-3 text-2xl"},"Select a namespace",-1),X=w({__name:"GroupView",setup(_){const s=k(!0),t=k([]);function m(){s.value=!0,o()}async function o(){const r=await C.fetchAll(i.namespace.value.id);t.value=Array.from(r)}return G(()=>{i.mustSelectNamespace||o()}),(r,e)=>(l(),c(B,null,[u("div",H,[J,n(i).mustSelectNamespace?y("",!0):(l(),c("button",{key:0,class:"ml-3 text-2xl border-none btn focus:outline-none txt-neutral",onClick:e[0]||(e[0]=a=>s.value=!s.value)},[s.value===!0?(l(),v(n(g),{key:0,icon:"fa6-solid:plus"})):(l(),v(n(g),{key:1,icon:"fa6-solid:minus"}))]))]),n(i).mustSelectNamespace?(l(),c("div",K,[L,d(D,{class:"mt-5",onSelectns:e[2]||(e[2]=a=>o())})])):(l(),c("div",{key:0,class:h([{"slide-y":!0,slideup:s.value===!0,slidedown:s.value===!1},"mb-8"])},[u("div",M,[T,s.value===!1?(l(),v(j,{key:0,class:"mt-5",onEnd:e[1]||(e[1]=a=>m())})):y("",!0)])],2)),n(i).mustSelectNamespace?y("",!0):(l(),v(U,{key:2,groups:t.value,onReload:e[3]||(e[3]=a=>o())},null,8,["groups"]))],64))}});export{X as default};
