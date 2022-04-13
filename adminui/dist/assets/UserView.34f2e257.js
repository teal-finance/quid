import{a as L,_ as j}from"./NamespaceSelector.b27b1aa6.js";import{d as V,o as a,c as d,a as u,F as q,e as I,f as k,b as v,n as g,a2 as C,j as F,i as b,w as _,l as G,t as H,u as t,s as N,v as $,k as O,S as U,m as T,y as z,I as R,p as K}from"./index.4dcf8b78.js";import{U as D}from"./user.7acceeab.js";import{s as A,A as S,a as M}from"./ActionButton.631ddfbc.js";import{_ as B}from"./SimpleBadge.4b0e2382.js";const J={key:0,class:"flex flex-col"},P={class:"mt-2 space-x-2"},Q=u("span",{class:"pr-2 font-bold"},"Groups:",-1),W={key:0},X={key:1,class:"txt-lighter"},Y={key:1},Z=V({props:{isLoading:{type:Boolean,default:!0},groups:{type:Array,required:!0},user:{type:Object,required:!0}},emits:["user-removed"],setup(f,{emit:y}){const r=f;function e(i){g.confirmDelete(`Remove the ${r.user.name} from group ${i.name}?`,()=>{C.removeUserFromGroup(r.user.id,i.id).then(()=>{g.done("User removed from group"),y("user-removed",i)})})}return(i,p)=>f.isLoading?(a(),d("div",Y,[v(L,{small:!0})])):(a(),d("div",J,[u("div",P,[Q,f.groups.length>0?(a(),d("span",W,[(a(!0),d(q,null,I(f.groups,n=>(a(),k(B,{key:n.id,class:"inline-block secondary cursor-pointer",text:n.name,onClick:w=>e(n)},null,8,["text","onClick"]))),128))])):(a(),d("span",X,"the user has no groups"))])]))}}),ee={key:0,class:"mt-5"},se=u("div",{class:"mt-5 flex flex-row pl-5 text-sm"},[u("div",{class:"mt-3 txt-lighter"},"Select a group")],-1),te={key:1,class:"mt-5 flex flex-row pl-5 text-sm"},ae=V({props:{user:{type:Object,required:!0},userGroups:{type:Array,required:!0}},emits:["addingroup"],setup(f,{emit:y}){const r=f,e=F({nsGroups:new Array}),i=b(0);async function p(){const w=await C.fetchAll($.namespace.value.id),o=r.userGroups.map(s=>s.id);e.nsGroups=[],w.forEach(s=>{o.includes(s.id)||e.nsGroups.push(s)}),i.value=1}function n(w){y("addingroup",w),i.value=0}return(w,o)=>i.value==1?(a(),d("div",ee,[(a(!0),d(q,null,I(t(e).nsGroups,s=>(a(),k(B,{key:s.id,class:"success cursor-pointer",onClick:l=>n(s)},{default:_(()=>[G(H(s.name),1)]),_:2},1032,["onClick"]))),128)),se])):i.value==0?(a(),d("div",te,[u("button",{class:"btn hover:secondary",onClick:o[0]||(o[0]=s=>p())},"Add user in group")])):N("",!0)}}),ne=G("Show groups"),re=G("Hide groups\xA0"),oe=G("Delete"),le={class:"p-3 pb-8 ml-5"},ie=V({props:{users:{type:Array,required:!0}},emits:["reload"],setup(f,{emit:y}){const r=f,e=b([]),i=b(0),p=F({userGroups:new Array}),n=b(!1);async function w(m,x){await C.addUserToGroup(x.id,m.id),p.userGroups.push(m)}async function o(m,x){p.userGroups=p.userGroups.filter(c=>{if(m.id!=c.id)return c})}function s(m){n.value=!0,C.fetchUserGroups(m).then(x=>{n.value=!1,p.userGroups.push(...x)}),e.value=r.users.filter(x=>x.id==m),i.value=m}function l(){p.userGroups=[],i.value=0,e.value=[]}function E(m){g.confirmDelete(`Delete the ${m.name} user?`,()=>{D.delete(m.id).then(()=>{g.done("User deleted"),y("reload")})})}return(m,x)=>(a(),d("div",null,[v(t(M),{value:f.users,class:"main-table p-datatable-sm",expandedRows:e.value,"onUpdate:expandedRows":x[1]||(x[1]=c=>e.value=c),"data-key":"id",removableSort:""},{expansion:_(c=>[u("div",le,[v(Z,{user:c.data,"is-loading":n.value,groups:t(p).userGroups,onUserRemoved:h=>o(h,c.data)},null,8,["user","is-loading","groups","onUserRemoved"]),v(ae,{user:c.data,onAddingroup:h=>w(h,c.data),"user-groups":t(p).userGroups},null,8,["user","onAddingroup","user-groups"])])]),default:_(()=>[v(t(A),{field:"id",header:"Id"}),v(t(A),{field:"name",header:"Name",sortable:!0}),v(t(A),{field:"actions"},{body:_(c=>[i.value!=c.data.id?(a(),k(S,{key:0,onClick:h=>s(c.data.id)},{default:_(()=>[ne]),_:2},1032,["onClick"])):(a(),k(S,{key:1,onClick:x[0]||(x[0]=h=>l())},{default:_(()=>[re]),_:1})),v(S,{type:"delete",class:"ml-2",onClick:h=>E(c.data)},{default:_(()=>[oe]),_:2},1032,["onClick"])]),_:1})]),_:1},8,["value","expandedRows"])]))}}),de={class:"flex flex-col space-y-5"},ue={class:"flex flex-row"},pe=["disabled"],ce={class:"inline-block ml-2"},ve=V({emits:["end"],setup(f,{emit:y}){const r=b(!1),e=F({name:{val:"",isValid:null,validator:o=>o.length>=3},pwd:{val:"",isValid:null,validator:o=>o.length>=5},pwdVerif:{val:"",isValid:null,validator:o=>o.length>=5}}),i=O(()=>e.name.isValid===!0&&e.pwd.isValid===!0&&e.pwdVerif.isValid===!0);function p(){if(e.pwd.val!==e.pwdVerif.val){g.warning("Password mismatch","Please verify the typing");return}n()}async function n(){try{await T.post("/admin/users/add",{name:e.name.val,password:e.pwd.val,namespace_id:$.namespace.value.id}),y("end"),g.done("User added")}catch(o){console.log(o),g.error("Error creating user")}}function w(){y("end")}return(o,s)=>(a(),d("div",de,[v(t(U),{value:t(e).name.val,"onUpdate:value":s[0]||(s[0]=l=>t(e).name.val=l),isvalid:t(e).name.isValid,"onUpdate:isvalid":s[1]||(s[1]=l=>t(e).name.isValid=l),validator:t(e).name.validator,"inline-label":"Name",required:"",autofocus:""},null,8,["value","isvalid","validator"]),v(t(U),{class:"mt-3",type:r.value?"text":"password",value:t(e).pwd.val,"onUpdate:value":s[2]||(s[2]=l=>t(e).pwd.val=l),isvalid:t(e).pwd.isValid,"onUpdate:isvalid":s[3]||(s[3]=l=>t(e).pwd.isValid=l),validator:t(e).pwd.validator,"inline-label":"Password",required:""},null,8,["type","value","isvalid","validator"]),v(t(U),{class:"mt-3",type:r.value?"text":"password",value:t(e).pwdVerif.val,"onUpdate:value":s[4]||(s[4]=l=>t(e).pwdVerif.val=l),isvalid:t(e).pwdVerif.isValid,"onUpdate:isvalid":s[5]||(s[5]=l=>t(e).pwdVerif.isValid=l),validator:t(e).pwdVerif.validator,"inline-label":"Password again",required:""},null,8,["type","value","isvalid","validator"]),u("div",ue,[u("button",{class:"w-20 mr-3 btn success",disabled:!t(i),onClick:s[6]||(s[6]=l=>p())},"Save",8,pe),u("button",{class:"w-20 btn warning",onClick:s[7]||(s[7]=l=>w())},"Cancel"),u("div",ce,[r.value?(a(),d("button",{key:0,class:"btn lighter",onClick:s[8]||(s[8]=l=>r.value=!1)},"Hide password")):(a(),d("button",{key:1,class:"btn lighter",onClick:s[9]||(s[9]=l=>r.value=!0)},"Show password"))])])]))}}),me={class:"text-3xl txt-primary dark:txt-light"},fe=G(" Users "),ye={class:"p-5 mt-3 border bord-lighter w-96"},we=u("div",{class:"text-xl"},"Add a user",-1),xe={key:1,class:"w-full"},_e=u("div",{class:"mt-3 text-2xl"},"Select a namespace",-1),Ve=V({setup(f){const y=b([]),r=b(!0);async function e(){y.value=await D.fetchAll($.namespace.value.id)}function i(){r.value=!0,e()}return z(()=>{$.mustSelectNamespace||e()}),(p,n)=>{const w=j;return a(),d(q,null,[u("div",me,[fe,t($).mustSelectNamespace?N("",!0):(a(),d("button",{key:0,class:"ml-3 text-2xl border-none btn focus:outline-none txt-neutral",onClick:n[0]||(n[0]=o=>r.value=!r.value)},[r.value===!0?(a(),k(t(R),{key:0,icon:"fa6-solid:plus"})):(a(),k(t(R),{key:1,icon:"fa6-solid:minus"}))]))]),t($).mustSelectNamespace?(a(),d("div",xe,[_e,v(w,{class:"mt-5",onSelectns:n[2]||(n[2]=o=>e())})])):(a(),d("div",{key:0,class:K([{"slide-y":!0,slideup:r.value===!0,slidedown:r.value===!1},"mb-8"])},[u("div",ye,[we,v(ve,{class:"mt-5",onEnd:n[1]||(n[1]=o=>i())})])],2)),t($).mustSelectNamespace?N("",!0):(a(),k(ie,{key:2,users:y.value,onReload:n[3]||(n[3]=o=>e())},null,8,["users"]))],64)}}});export{Ve as default};
