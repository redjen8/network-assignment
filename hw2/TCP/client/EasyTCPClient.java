package hw2.TCP.client;

import java.net.*;
import java.util.Scanner;
import java.io.*;

public class EasyTCPClient {

    private static final String SERVER_NAME = "localhost";
    private static final int SERVER_PORT = 22848;
    
    public static void main(String[] args) {
        Socket socket = null;

        try {
            socket = new Socket();
            socket.connect(new InetSocketAddress(SERVER_NAME, SERVER_PORT));
            byte[] bytes = null;
            String message = null;
            OutputStream os = socket.getOutputStream();

            Scanner sc = new Scanner(System.in);

            
            message = "ASK_IP_PORT";
            bytes = message.getBytes("UTF-8");
            os.write(bytes);
            os.flush();

            InputStream is = socket.getInputStream();
            bytes = new byte[1024];
            int readByteCount = is.read(bytes);
            message = new String(bytes, 0, readByteCount, "UTF-8");
            System.out.println("message : " + message);
            os.close();
            is.close();
            socket.close();
        }
        catch (IOException e) {
            e.printStackTrace();
        }
        if (!socket.isClosed()) {
            try {
                socket.close();
            }
            catch (IOException e) {
                e.printStackTrace();
            }
        }

    }
}