import{m as o}from"./index.f6d5cb7b.js";class c{constructor(s){this.id=s.id,this.name=s.username}toTableRow(){const s=Object.create(this);return s.actions=[],s}static async search(s,r){const e="/admin/users/search",t=new Array;try{const a={namespace_id:s,username:r};(await o.post(e,a)).users.forEach(i=>t.push(new c(i)))}catch(a){throw console.log("Err",a),a}return t}static async fetchAll(s){const r="/admin/users/nsall",e=new Array;try{const t={namespace_id:s};(await o.post(r,t)).forEach(n=>e.push(new c(n).toTableRow()))}catch(t){throw console.log("Err",t),t}return e}static async delete(s){await o.post("/admin/users/delete",{id:s})}}export{c as U};
