import { Migrations } from './index' 

export default class Workspaces {
    
    /**
     * loadCurrenWorkspace
     * @param initialState
     * @returns {}
     */
    static loadCurrentWorkspace(initialState){
        const current_workspace = Workspaces.loadFromLocalStorage("current_workspace");

        if(current_workspace){
            const { state, version } = current_workspace;

            const compatible = Workspaces.checkVersions(version);
            if(!compatible){
                return Workspaces.migrate(version, state);
            }

            return Workspaces.actual(version, state);
        }

        return initialState;
    }


    /**
     * persistCurrentWorkspace
     * @param Object state
     * @return Object
     */
    static persistCurrentWorkspace(state){
        try {

            const persistedState = Object.assign({}, {state: state}, {
                version: Migrations.getCurrentVersion(),
                state: {
                    ...state,
                    routing: {},
                    entries:{
                        ...state.entries,
                        node: [],
                        items: [],
                        nodes: [],
                        links: [],
                        searches: [],
                    }
                }
            });

            const string = JSON.stringify(persistedState);
            localStorage.setItem("current_workspace", string);
        } catch(e){
            console.warn('Unable to persist state to localStorage:', e);
        }
    }


    /**
     * migrate
     * @param int version
     * @param Object state
     * @return Object
     */
    static migrate(version, state){
        const migrations = Migrations.getMigrationsToApply(version, Migrations.getCurrentVersion());
        console.info(`Applying ${migrations.length} migrations`);         
        return Migrations.applyAllMigrations(version, state, migrations);
    }


    /**
     * actual
     * @param int version
     * @param Object state
     * @return Object
     */
    static actual(version, state){
        console.info(`Found compatible state version: ${version} in localStorage, recklessly persisting to application state`);
        return state;
    }

        
    /**
     * checkVersions
     * @param int version
     * @return bool
     */
    static checkVersions(version){
        return parseInt(version) == Migrations.getCurrentVersion();
    }


    /**
     * loadFromLocalStorage
     * @param key
     * @return {Object|bool}
     */
    static loadFromLocalStorage(key){
        try {
            return JSON.parse(localStorage.getItem(key));
        }catch(e){
            return false
        }
    }
}
