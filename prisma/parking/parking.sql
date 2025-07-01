CREATE TABLE lot (
                     lot_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                     prime_location_name TEXT NOT NULL,
                     price INTEGER NOT NULL,
                     address TEXT NOT NULL,
                     pincode TEXT NOT NULL,
                     spots INTEGER NOT NULL
);


CREATE TABLE IF NOT EXISTS spot (
                                    parking_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                    lot_id INT NOT NULL,
                                    status BOOLEAN NOT NULL,
                                    address TEXT NOT NULL,
                                    CONSTRAINT fk_lot
                                    FOREIGN KEY (lot_id)
    REFERENCES lot(lot_id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
    );

CREATE TABLE IF NOT EXISTS "reserved_spot" (
                                               "id" TEXT PRIMARY KEY NOT NULL,
                                               "spot_id" INT NOT NULL,
                                               "user_id" TEXT NOT NULL,
                                                "vehicle_no" TEXT NOT NULL,
                                                "location"   TEXT NOT NULL,
                                               "parking" TIMESTAMP NOT NULL,
                                               "leaving" TIMESTAMP NOT NULL,
                                               "parking_cost" INTEGER NOT NULL,
                                                "status"    BOOLEAN NOT NULL,

                                               CONSTRAINT fk_user
                                               FOREIGN KEY ("user_id")
    REFERENCES "user_data"("user_name")
    ON DELETE CASCADE
    ON UPDATE CASCADE,

    CONSTRAINT fk_spot
    FOREIGN KEY ("spot_id")
    REFERENCES "spot"("parking_id")
    ON DELETE CASCADE
    ON UPDATE CASCADE
    );
