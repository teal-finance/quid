import { User } from "@snowind/state";
import { useApi } from '@snowind/api';
import { libName } from "./conf";
import { reactive } from "vue";

const user = new User();
const api = useApi({ serverUrl: import.meta.env.MODE === 'development' ? '' : `/${libName.toLowerCase()}` });
const state = reactive({
  apidocs: new Array<string>(),
  server: new Array<string>(),
  client: new Array<string>(),
  examples: new Array<string>(),
})

function fetchIndexes() {
  api.get<Array<string>>("/apidoc/index.json").then((res) => {
    state.apidocs = typeof res == "string" ? JSON.parse(res) : res
  });
  api.get<Array<string>>("/server/index.json").then(res =>
    state.server = typeof res == "string" ? JSON.parse(res) : res
  );
  api.get<Array<string>>("/client/index.json").then((res) => {
    state.client = typeof res == "string" ? JSON.parse(res) : res
  });
  api.get<Array<string>>("/examples/index.json").then(res =>
    state.examples = typeof res == "string" ? JSON.parse(res) : res
  );
}

function initState() {
  fetchIndexes()
}

export { user, api, state, initState }