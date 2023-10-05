package com.samaya.gotictactoe;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;

public class Util {

    public static byte[] readFully(InputStream is) throws IOException {
        byte retval[] = null;

        if(is != null){
            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            byte packet[] = new byte[1024];

            while(is.read(packet) > 0 ){
                baos.write(packet);
            }
            retval = baos.toByteArray();
            baos.close();
        }
        return retval;
    }
}
