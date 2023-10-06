package com.samaya.gotictactoe;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Date;

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


    public static String formatDate(Date date){
        String retval = "";
        if(date != null){
            retval = String.format("%td %tb %tY", date, date, date);
        }
        return  retval;
    }

    public static String formatTime(long duration){
        long hour = duration / 60;
        long min = duration % 60;

        return String.format("%d : %d", hour, min);
    }
}
