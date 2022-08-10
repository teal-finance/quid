import{a as L,_ as T}from"./NamespaceSelector-988775c6.js";import{h as x,z as I,d as G,o as a,c as d,a as c,F as q,e as B,f as k,b as m,A as h,a5 as U,x as F,j as b,w as $,B as V,t as H,H as N,u as v,y as O,S as C,J as z,I as D,p as J}from"./index-6b0967df.js";import{s as K,a as S,A as R}from"./ActionButton-36ce7124.js";import{_ as E}from"./SimpleBadge-0cbc5bef.js";class A{constructor(o){this.id=o.id,this.name=o.username}toTableRow(){const o=Object.create(this);return o.actions=[],o}static async fetchAll(o){const t=x.adminUrl+"/users/nsall",e=new Array;try{const n={namespace_id:o};(await I.post(t,n)).forEach(r=>e.push(new A(r).toTableRow()))}catch(n){throw console.log("Err",n),n}return e}static async delete(o){const t=x.adminUrl+"/users/delete";await I.post(t,{id:o,namespace_id:x.namespace.value.id})}}const M={key:0,class:"flex flex-col"},P={class:"mt-2 space-x-2"},Q=c("span",{class:"pr-2 font-bold"},"Groups:",-1),W={key:0},X={key:1,class:"txt-lighter"},Y={key:1},Z=G({__name:"UserGroupsInfo",props:{isLoading:{type:Boolean,default:!0},groups:{type:Array,required:!0},user:{type:Object,required:!0}},emits:["user-removed"],setup(w,{emit:o}){const t=w;function e(n){h.confirmDelete(`Remove the ${t.user.name} from group ${n.name}?`,()=>{U.removeUserFromGroup(t.user.id,n.id).then(()=>{h.done("User removed from group"),o("user-removed",n)})})}return(n,u)=>w.isLoading?(a(),d("div",Y,[m(L,{small:!0})])):(a(),d("div",M,[c("div",P,[Q,w.groups.length>0?(a(),d("span",W,[(a(!0),d(q,null,B(w.groups,r=>(a(),k(E,{key:r.id,class:"inline-block secondary cursor-pointer",text:r.name,onClick:y=>e(r)},null,8,["text","onClick"]))),128))])):(a(),d("span",X,"the user has no groups"))])]))}}),ee={key:0,class:"mt-5"},se=c("div",{class:"mt-5 flex flex-row pl-5 text-sm"},[c("div",{class:"mt-3 txt-lighter"},"Select a group")],-1),te={key:1,class:"mt-5 flex flex-row pl-5 text-sm"},ae=G({__name:"AddUserIntoGroup",props:{user:{type:Object,required:!0},userGroups:{type:Array,required:!0}},emits:["addingroup"],setup(w,{emit:o}){const t=w,e=F({nsGroups:new Array}),n=b(0);async function u(){const y=await U.fetchAll(x.namespace.value.id),l=t.userGroups.map(s=>s.id);e.nsGroups=[],y.forEach(s=>{l.includes(s.id)||e.nsGroups.push(s)}),n.value=1}function r(y){o("addingroup",y),n.value=0}return(y,l)=>n.value==1?(a(),d("div",ee,[(a(!0),d(q,null,B(e.nsGroups,s=>(a(),k(E,{key:s.id,class:"success cursor-pointer",onClick:i=>r(s)},{default:$(()=>[V(H(s.name),1)]),_:2},1032,["onClick"]))),128)),se])):n.value==0?(a(),d("div",te,[c("button",{class:"btn hover:secondary",onClick:l[0]||(l[0]=s=>u())},"Add user in group")])):N("",!0)}}),ne=V("Show groups"),re=V("Hide groups\xA0"),oe=V("Delete"),le={class:"p-3 pb-8 ml-5"},ie=G({__name:"UserDatatable",props:{users:{type:Array,required:!0}},emits:["reload"],setup(w,{emit:o}){const t=w,e=b([]),n=b(0),u=F({userGroups:new Array}),r=b(!1);async function y(f,_){await U.addUserToGroup(_.id,f.id),u.userGroups.push(f)}async function l(f,_){u.userGroups=u.userGroups.filter(p=>{if(f.id!=p.id)return p})}function s(f){r.value=!0,U.fetchUserGroups(f).then(_=>{r.value=!1,u.userGroups.push(..._)}),e.value=t.users.filter(_=>_.id==f),n.value=f}function i(){u.userGroups=[],n.value=0,e.value=[]}function j(f){h.confirmDelete(`Delete the ${f.name} user?`,()=>{A.delete(f.id).then(()=>{h.done("User deleted"),o("reload")})})}return(f,_)=>(a(),d("div",null,[m(v(K),{value:w.users,class:"main-table p-datatable-sm",expandedRows:e.value,"onUpdate:expandedRows":_[1]||(_[1]=p=>e.value=p),"data-key":"id",removableSort:""},{expansion:$(p=>[c("div",le,[m(Z,{user:p.data,"is-loading":r.value,groups:u.userGroups,onUserRemoved:g=>l(g,p.data)},null,8,["user","is-loading","groups","onUserRemoved"]),m(ae,{user:p.data,onAddingroup:g=>y(g,p.data),"user-groups":u.userGroups},null,8,["user","onAddingroup","user-groups"])])]),default:$(()=>[m(v(S),{field:"id",header:"Id"}),m(v(S),{field:"name",header:"Name",sortable:!0}),m(v(S),{field:"actions"},{body:$(p=>[n.value!=p.data.id?(a(),k(R,{key:0,onClick:g=>s(p.data.id)},{default:$(()=>[ne]),_:2},1032,["onClick"])):(a(),k(R,{key:1,onClick:_[0]||(_[0]=g=>i())},{default:$(()=>[re]),_:1})),m(R,{type:"delete",class:"ml-2",onClick:g=>j(p.data)},{default:$(()=>[oe]),_:2},1032,["onClick"])]),_:1})]),_:1},8,["value","expandedRows"])]))}}),de={class:"flex flex-col space-y-5"},ue={class:"flex flex-row"},ce=["disabled"],pe={class:"inline-block ml-2"},me=G({__name:"AddUser",emits:["end"],setup(w,{emit:o}){const t=b(!1),e=F({name:{val:"",isValid:null,validator:l=>l.length>=3},pwd:{val:"",isValid:null,validator:l=>l.length>=5},pwdVerif:{val:"",isValid:null,validator:l=>l.length>=5}}),n=O(()=>e.name.isValid===!0&&e.pwd.isValid===!0&&e.pwdVerif.isValid===!0);function u(){if(e.pwd.val!==e.pwdVerif.val){h.warning("Password mismatch","Please verify the typing");return}r()}async function r(){try{await I.post(x.adminUrl+"/users/add",{name:e.name.val,password:e.pwd.val,namespace_id:x.namespace.value.id}),o("end"),h.done("User added")}catch(l){console.log(l),h.error("Error creating user")}}function y(){o("end")}return(l,s)=>(a(),d("div",de,[m(v(C),{value:e.name.val,"onUpdate:value":s[0]||(s[0]=i=>e.name.val=i),isvalid:e.name.isValid,"onUpdate:isvalid":s[1]||(s[1]=i=>e.name.isValid=i),validator:e.name.validator,"inline-label":"Name",required:"",autofocus:""},null,8,["value","isvalid","validator"]),m(v(C),{class:"mt-3",type:t.value?"text":"password",value:e.pwd.val,"onUpdate:value":s[2]||(s[2]=i=>e.pwd.val=i),isvalid:e.pwd.isValid,"onUpdate:isvalid":s[3]||(s[3]=i=>e.pwd.isValid=i),validator:e.pwd.validator,"inline-label":"Password",required:""},null,8,["type","value","isvalid","validator"]),m(v(C),{class:"mt-3",type:t.value?"text":"password",value:e.pwdVerif.val,"onUpdate:value":s[4]||(s[4]=i=>e.pwdVerif.val=i),isvalid:e.pwdVerif.isValid,"onUpdate:isvalid":s[5]||(s[5]=i=>e.pwdVerif.isValid=i),validator:e.pwdVerif.validator,"inline-label":"Password again",required:""},null,8,["type","value","isvalid","validator"]),c("div",ue,[c("button",{class:"w-20 mr-3 btn success",disabled:!v(n),onClick:s[6]||(s[6]=i=>u())},"Save",8,ce),c("button",{class:"w-20 btn warning",onClick:s[7]||(s[7]=i=>y())},"Cancel"),c("div",pe,[t.value?(a(),d("button",{key:0,class:"btn lighter",onClick:s[8]||(s[8]=i=>t.value=!1)},"Hide password")):(a(),d("button",{key:1,class:"btn lighter",onClick:s[9]||(s[9]=i=>t.value=!0)},"Show password"))])])]))}}),ve={class:"text-3xl txt-primary dark:txt-light"},fe=V(" Users "),we={class:"p-5 mt-3 border bord-lighter w-96"},ye=c("div",{class:"text-xl"},"Add a user",-1),_e={key:1,class:"w-full"},xe=c("div",{class:"mt-3 text-2xl"},"Select a namespace",-1),ge=G({__name:"UserView",setup(w){const o=b([]),t=b(!0);async function e(){o.value=await A.fetchAll(x.namespace.value.id)}function n(){t.value=!0,e()}return z(()=>{x.mustSelectNamespace||e()}),(u,r)=>{const y=T;return a(),d(q,null,[c("div",ve,[fe,v(x).mustSelectNamespace?N("",!0):(a(),d("button",{key:0,class:"ml-3 text-2xl border-none btn focus:outline-none txt-neutral",onClick:r[0]||(r[0]=l=>t.value=!t.value)},[t.value===!0?(a(),k(v(D),{key:0,icon:"fa6-solid:plus"})):(a(),k(v(D),{key:1,icon:"fa6-solid:minus"}))]))]),v(x).mustSelectNamespace?(a(),d("div",_e,[xe,m(y,{class:"mt-5",onSelectns:r[2]||(r[2]=l=>e())})])):(a(),d("div",{key:0,class:J([{"slide-y":!0,slideup:t.value===!0,slidedown:t.value===!1},"mb-8"])},[c("div",we,[ye,m(me,{class:"mt-5",onEnd:r[1]||(r[1]=l=>n())})])],2)),v(x).mustSelectNamespace?N("",!0):(a(),k(ie,{key:2,users:o.value,onReload:r[3]||(r[3]=l=>e())},null,8,["users"]))],64)}}});export{ge as default};