import{C as $,_ as u,c as h}from"./index.3f7e8321.js";import{d as c,o as t,b as s,a as i,f as n,u as v,I as y,t as x,r as b,F as w,e as B,c as d,w as F,p as V}from"./vendor.ef2cb0e3.js";const j={class:"p-4 text-lg text-center bg-transparent border shadow dark:bg-gray-500"},q={class:"flex justify-center w-full"},I={class:"mt-3 mb-3 dark:font-bold"},_=c({props:{card:{type:$,required:!0}},setup(e){return(a,l)=>(t(),s("div",j,[i("div",q,[n(v(y),{icon:e.card.icon,class:"mt-2 text-3xl txt-neutral"},null,8,["icon"])]),i("div",I,x(e.card.title),1)]))}}),N=c({components:{SimpleCard:_},props:{cards:{type:Object,required:!0},onCtrlClickCard:{type:Function,default:()=>null}}}),S={class:"grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6"},D=["onClick"];function E(e,a,l,m,f,g){const o=_,C=b("router-link");return t(),s("div",S,[(t(!0),s(w,null,B(e.cards,(r,k)=>(t(),s("div",{key:k},[r.url?(t(),d(C,{key:0,to:r.url},{default:F(()=>[n(o,{card:r},null,8,["card"])]),_:2},1032,["to"])):e.onCtrlClickCard!==null?(t(),s("div",{key:1,onClick:V(M=>e.onCtrlClickCard(r),["ctrl"])},[n(o,{card:r},null,8,["card"])],8,D)):(t(),d(o,{key:2,card:r},null,8,["card"]))]))),128))])}var p=u(N,[["render",E]]);const G=c({components:{CardsGrid:p},setup(){return{categories:h}}}),H={class:"p-5"};function L(e,a,l,m,f,g){const o=p;return t(),s("div",H,[n(o,{cards:e.categories},null,8,["cards"])])}var A=u(G,[["render",L]]);export{A as default};
