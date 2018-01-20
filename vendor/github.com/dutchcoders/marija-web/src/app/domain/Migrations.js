const migrations = {
    2: (state) => migration(`This doesn't do anything, it's here as an example, for when we want to implement our first migration.`, (inner_state) => {
            return inner_state;
    }, state),
}

function migration(message, callback, state){
    console.debug(message);

    const oldState = Object.assign({},state);
    return callback(oldState);
}

export default class Migrations {

    /**
     * loadCurrenWorkspace
     * @param initialState
     * @param states
     * @returns int
     */
    static getCurrentVersion(){
        return 3;
    }


    /**
     * getMigrationsToApply
     * @param version
     */
    static getMigrationsToApply(version, targetVersion){
        return Object.keys(migrations).filter((item) => {
            return (item > version && item <= targetVersion);
        }).map((key) => {
            return migrations[key];
        });
    }


    /**
     * applyAllMigrations
     * @param version
     * @param state
     * @param migrations
     */
    static applyAllMigrations(version, state, migrations){
        return migrations.reduce((endState, migration) => {
            const nextVersion = version + (migrations.indexOf(migration) + 1);

            console.info(`Applying version ${nextVersion} to version ${(nextVersion - 1)} `);
            return migration(endState);
        }, state)
    }

}
